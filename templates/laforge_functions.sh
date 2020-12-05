#!/bin/sh -
#======================================================================================================================
# vim: softtabstop=2 shiftwidth=2 expandtab fenc=utf-8 spell spelllang=en cc=120
#======================================================================================================================
#
#          FILE: laforge_functions.sh
#
#   DESCRIPTION: Laforge function library.
#
#          BUGS: https://github.com/gen0cide/laforge/issues
#
#======================================================================================================================
set -o nounset # Treat unset variables as an error

# Bootstrap script truth values
export LAFORGE_BS_TRUE=1
export LAFORGE_BS_FALSE=0
export LAFORGE_DIR="/root/.laforge"

export DEBIAN_FRONTEND=noninteractive

# Default sleep time used when waiting for daemons to start, restart and checking for these running
export LAFORGE_DEFAULT_SLEEP=3
export LAFORGE_CURL_ARGS="--insecure"
export LAFORGE_FETCH_ARGS="--no-verify-peer"
export LAFORGE_GPG_ARGS="--keyserver-options no-check-cert"
export LAFORGE_WGET_ARGS="--no-check-certificate"
export LAFORGE_DEBUG=1

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  __detect_color_support
#   DESCRIPTION:  Try to detect color support.
#----------------------------------------------------------------------------------------------------------------------
export LAFORGE_COLORS=${BS_COLORS:-$(tput colors 2>/dev/null || echo 0)}
laforge_detect_color_support() {
	# shellcheck disable=SC2181
	if [ $? -eq 0 ] && [ "$LAFORGE_COLORS" -gt 2 ]; then
		export LAFORGE_RC='\033[1;31m'
		export LAFORGE_GC='\033[1;32m'
		export LAFORGE_BC='\033[1;34m'
		export LAFORGE_YC='\033[1;33m'
		export LAFORGE_EC='\033[0m'
	else
		export LAFORGE_RC=""
		export LAFORGE_GC=""
		export LAFORGE_BC=""
		export LAFORGE_YC=""
		export LAFORGE_EC=""
	fi
}

laforge_detect_color_support

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  echoerr
#   DESCRIPTION:  Echo errors to stderr.
#----------------------------------------------------------------------------------------------------------------------
laforge_error() {
	printf "${LAFORGE_RC} * ERROR${LAFORGE_EC}: %s\\n" "$@" 1>&2
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  echoinfo
#   DESCRIPTION:  Echo information to stdout.
#----------------------------------------------------------------------------------------------------------------------
laforge_info() {
	printf "${LAFORGE_GC} *  INFO${LAFORGE_EC}: %s\\n" "$@"
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  echowarn
#   DESCRIPTION:  Echo warning information to stdout.
#----------------------------------------------------------------------------------------------------------------------
laforge_warn() {
	printf "${LAFORGE_YC} *  WARN${LAFORGE_EC}: %s\\n" "$@"
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  echodebug
#   DESCRIPTION:  Echo debug information to stdout.
#----------------------------------------------------------------------------------------------------------------------
laforge_debug() {
	if [ "$LAFORGE_DEBUG" -eq $LAFORGE_BS_TRUE ]; then
		printf "${LAFORGE_BC} * DEBUG${LAFORGE_EC}: %s\\n" "$@"
	fi
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_check_command_exists
#   DESCRIPTION:  Check if a command exists.
#----------------------------------------------------------------------------------------------------------------------
laforge_check_command_exists() {
	command -v "$1" >/dev/null 2>&1
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#         NAME:  laforge_fetch_url
#  DESCRIPTION:  Retrieves a URL and writes it to a given path
#----------------------------------------------------------------------------------------------------------------------
laforge_fetch_url() {
	# shellcheck disable=SC2086
	curl $LAFORGE_CURL_ARGS -L -s -o "$1" "$2" >/dev/null 2>&1 ||
		wget $LAFORGE_WGET_ARGS -q -O "$1" "$2" >/dev/null 2>&1 ||
		fetch $LAFORGE_FETCH_ARGS -q -o "$1" "$2" >/dev/null 2>&1 || # FreeBSD
		fetch -q -o "$1" "$2" >/dev/null 2>&1 || # Pre FreeBSD 10
		ftp -o "$1" "$2" >/dev/null 2>&1 # OpenBSD
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_gather_hardware_info
#   DESCRIPTION:  Discover hardware information
#----------------------------------------------------------------------------------------------------------------------
laforge_gather_hardware_info() {
	if [ -f /proc/cpuinfo ]; then
		export LAFORGE_CPU_VENDOR_ID=$(awk '/vendor_id|Processor/ {sub(/-.*$/,"",$3); print $3; exit}' /proc/cpuinfo)
	elif [ -f /usr/bin/kstat ]; then
		# SmartOS.
		# Solaris!?
		# This has only been tested for a GenuineIntel CPU
		export LAFORGE_CPU_VENDOR_ID=$(/usr/bin/kstat -p cpu_info:0:cpu_info0:vendor_id | awk '{print $2}')
	else
		export LAFORGE_CPU_VENDOR_ID=$(sysctl -n hw.model)
	fi
	export LAFORGE_CPU_VENDOR_ID_L=$(echo "$LAFORGE_CPU_VENDOR_ID" | tr '[:upper:]' '[:lower:]')
	export LAFORGE_CPU_ARCH=$(uname -m 2>/dev/null || uname -p 2>/dev/null || echo "unknown")
	export LAFORGE_CPU_ARCH_L=$(echo "$LAFORGE_CPU_ARCH" | tr '[:upper:]' '[:lower:]')
}
laforge_gather_hardware_info

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_gather_os_info
#   DESCRIPTION:  Discover operating system information
#----------------------------------------------------------------------------------------------------------------------
laforge_gather_os_info() {
	export LAFORGE_OS_NAME=$(uname -s 2>/dev/null)
	export LAFORGE_OS_NAME_L=$(echo "$LAFORGE_OS_NAME" | tr '[:upper:]' '[:lower:]')
	export LAFORGE_OS_VERSION=$(uname -r)
	# shellcheck disable=SC2034 SC2155
	export LAFORGE_OS_VERSION_L=$(echo "$LAFORGE_OS_VERSION" | tr '[:upper:]' '[:lower:]')
}
laforge_gather_os_info

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_wait_for_apt
#   DESCRIPTION:  Check if any apt, apt-get, aptitude, or dpkg processes are running before
#                 calling these again. This is useful when these process calls are part of
#                 a boot process, such as on AWS AMIs. This func will wait until the boot
#                 process is finished so the script doesn't exit on a locked proc.
#----------------------------------------------------------------------------------------------------------------------
laforge_wait_for_apt() {
	laforge_debug "Checking if apt process is currently running."

	# Timeout set at 15 minutes
	WAIT_TIMEOUT=900

	while ps -C apt,apt-get,aptitude,dpkg >/dev/null; do
		sleep 1
		WAIT_TIMEOUT=$((WAIT_TIMEOUT - 1))

		# If timeout reaches 0, abort.
		if [ "$WAIT_TIMEOUT" -eq 0 ]; then
			laforge_error "Apt, apt-get, aptitude, or dpkg process is taking too long."
			laforge_error "Bootstrap script cannot proceed. Aborting."
			return 1
		fi
	done

	laforge_debug "No apt processes are currently running."
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_apt_install
#   DESCRIPTION:  (DRY) apt-get install with noinput options
#    PARAMETERS:  packages
#----------------------------------------------------------------------------------------------------------------------
function laforge_apt_install() {
	laforge_wait_for_apt
	apt-get install -y -o DPkg::Options::=--force-confold "${@}"
	return $?
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_apt_get_upgrade
#   DESCRIPTION:  (DRY) apt-get upgrade with noinput options
#----------------------------------------------------------------------------------------------------------------------
function laforge_apt_get_upgrade() {
	laforge_wait_for_apt
	apt-get upgrade -y -o DPkg::Options::=--force-confold
	return $?
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_apt_update
#   DESCRIPTION:  apt-get update with suppressed output
#    PARAMETERS:  none
#----------------------------------------------------------------------------------------------------------------------
laforge_apt_update() {
	laforge_wait_for_apt

	laforge_info "Updating apt repositories."
	apt-get -qq update || return 1
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_apt_key_fetch
#   DESCRIPTION:  Download and import GPG public key for "apt-secure"
#    PARAMETERS:  url
#----------------------------------------------------------------------------------------------------------------------
laforge_apt_key_fetch() {
	laforge_wait_for_apt
	url=$1

	# shellcheck disable=SC2086
	apt-key adv ${LAFORGE_GPG_ARGS} --fetch-keys "$url"
	return $?
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_disable_output
#   DESCRIPTION:  Creates a log file and forces STDERR and STDOUT to it
#    PARAMETERS:  none
#----------------------------------------------------------------------------------------------------------------------
function laforge_disable_output() {
	mkdir -p "${LAFORGE_DIR}"
	# Define our logging file and pipe paths
	LOGFILE="${LAFORGE_DIR}/$(echo "$1" | sed s/.sh/.log/g)"
	LOGPIPE="/tmp/$(echo "$1" | sed s/.sh/.logpipe/g)"
	# Ensure no residual pipe exists
	rm "$LOGPIPE" 2>/dev/null

	# Create our logging pipe
	# On FreeBSD we have to use mkfifo instead of mknod
	if ! (mknod "$LOGPIPE" p >/dev/null 2>&1 || mkfifo "$LOGPIPE" >/dev/null 2>&1); then
		echoerror "Failed to create the named pipe required to log"
		exit 1
	fi

	# What ever is written to the logpipe gets written to the logfile
	tee <"$LOGPIPE" "$LOGFILE" &

	# Close STDOUT, reopen it directing it to the logpipe
	exec 1>&-
	exec 1>"$LOGPIPE"
	# Close STDERR, reopen it directing it to the logpipe
	exec 2>&-
	exec 2>"$LOGPIPE"
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_check_services_systemd
#   DESCRIPTION:  Return 0 or 1 in case the service is enabled or not
#    PARAMETERS:  servicename
#----------------------------------------------------------------------------------------------------------------------
laforge_check_services_systemd() {
	if [ $# -eq 0 ]; then
		laforge_error "You need to pass a service name to check!"
		exit 1
	elif [ $# -ne 1 ]; then
		laforge_error "You need to pass a service name to check as the single argument to the function"
	fi

	servicename=$1
	laforge_debug "Checking if service ${servicename} is enabled"

	if [ "$(systemctl is-enabled "${servicename}")" = "enabled" ]; then
		laforge_debug "Service ${servicename} is enabled"
		return 0
	else
		laforge_debug "Service ${servicename} is NOT enabled"
		return 1
	fi
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_enable_universe_repository
#   DESCRIPTION:  Enable the universe repository if it is not already enabled
#    PARAMETERS:  none
#----------------------------------------------------------------------------------------------------------------------
laforge_enable_universe_repository() {
	if [ "$(grep -R universe /etc/apt/sources.list /etc/apt/sources.list.d/ | grep -v '#')" != "" ]; then
		# The universe repository is already enabled
		return 0
	fi

	laforge_debug "Enabling the universe repository"
	add-apt-repository -y "deb http://archive.ubuntu.com/ubuntu $(lsb_release -sc) universe" || return 1
	return 0
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_ubuntu_prep
#   DESCRIPTION:  Enable the universe repository and do an apt update
#    PARAMETERS:  none
#----------------------------------------------------------------------------------------------------------------------
laforge_ubuntu_prep() {
	# Install add-apt-repository
	if ! laforge_check_command_exists add-apt-repository; then
		laforge_apt_install software-properties-common || return 1
	fi

	laforge_enable_universe_repository || return 1

	laforge_apt_update
}

#---  FUNCTION  -------------------------------------------------------------------------------------------------------
#          NAME:  laforge_wait_for_port
#   DESCRIPTION:  Wait for up to five minutes for a port to become available
#    PARAMETERS:  none
#----------------------------------------------------------------------------------------------------------------------
laforge_wait_for_port() {
	if [ $# -eq 0 ]; then
		laforge_error "You need to pass a port to monitor!"
		exit 1
	elif [ $# -ne 1 ]; then
		laforge_error "You need to pass a single argument (a port) to the function"
	fi
	port=$1

	# Timeout set at 5 minutes
	WAIT_TIMEOUT=300

	while !(nc -z localhost ${port}) >/dev/null; do
		sleep 1
		let WAIT_TIMEOUT-=1

		# If timeout reaches 0, abort.
		if [ "$WAIT_TIMEOUT" -eq 0 ]; then
			laforge_error "Waiting for port ${port} to open has taken too long."
			laforge_error "Script cannot proceed. Aborting."
			return 1
		fi
	done
}

#----------------------------------------------------------------------------------------------------------------------
#  Handle command line arguments
#----------------------------------------------------------------------------------------------------------------------
