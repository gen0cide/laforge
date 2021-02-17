// Code generated by entc, DO NOT EDIT.

package ent

import "context"

func (as *AgentStatus) AgentStatusToTag(ctx context.Context) ([]*Tag, error) {
	result, err := as.Edges.AgentStatusToTagOrErr()
	if IsNotLoaded(err) {
		result, err = as.QueryAgentStatusToTag().All(ctx)
	}
	return result, err
}

func (as *AgentStatus) AgentStatusToProvisionedHost(ctx context.Context) ([]*ProvisionedHost, error) {
	result, err := as.Edges.AgentStatusToProvisionedHostOrErr()
	if IsNotLoaded(err) {
		result, err = as.QueryAgentStatusToProvisionedHost().All(ctx)
	}
	return result, err
}

func (b *Build) BuildToUser(ctx context.Context) ([]*User, error) {
	result, err := b.Edges.BuildToUserOrErr()
	if IsNotLoaded(err) {
		result, err = b.QueryBuildToUser().All(ctx)
	}
	return result, err
}

func (b *Build) BuildToTag(ctx context.Context) ([]*Tag, error) {
	result, err := b.Edges.BuildToTagOrErr()
	if IsNotLoaded(err) {
		result, err = b.QueryBuildToTag().All(ctx)
	}
	return result, err
}

func (b *Build) BuildToProvisionedNetwork(ctx context.Context) ([]*ProvisionedNetwork, error) {
	result, err := b.Edges.BuildToProvisionedNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = b.QueryBuildToProvisionedNetwork().All(ctx)
	}
	return result, err
}

func (b *Build) BuildToTeam(ctx context.Context) ([]*Team, error) {
	result, err := b.Edges.BuildToTeamOrErr()
	if IsNotLoaded(err) {
		result, err = b.QueryBuildToTeam().All(ctx)
	}
	return result, err
}

func (b *Build) BuildToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := b.Edges.BuildToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = b.QueryBuildToEnvironment().All(ctx)
	}
	return result, err
}

func (c *Command) CommandToUser(ctx context.Context) ([]*User, error) {
	result, err := c.Edges.CommandToUserOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryCommandToUser().All(ctx)
	}
	return result, err
}

func (c *Command) CommandToTag(ctx context.Context) ([]*Tag, error) {
	result, err := c.Edges.CommandToTagOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryCommandToTag().All(ctx)
	}
	return result, err
}

func (c *Competition) CompetitionToTag(ctx context.Context) ([]*Tag, error) {
	result, err := c.Edges.CompetitionToTagOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryCompetitionToTag().All(ctx)
	}
	return result, err
}

func (c *Competition) CompetitionToDNS(ctx context.Context) ([]*DNS, error) {
	result, err := c.Edges.CompetitionToDNSOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryCompetitionToDNS().All(ctx)
	}
	return result, err
}

func (c *Competition) CompetitionToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := c.Edges.CompetitionToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = c.QueryCompetitionToEnvironment().All(ctx)
	}
	return result, err
}

func (d *DNS) DNSToTag(ctx context.Context) ([]*Tag, error) {
	result, err := d.Edges.DNSToTagOrErr()
	if IsNotLoaded(err) {
		result, err = d.QueryDNSToTag().All(ctx)
	}
	return result, err
}

func (dr *DNSRecord) DNSRecordToTag(ctx context.Context) ([]*Tag, error) {
	result, err := dr.Edges.DNSRecordToTagOrErr()
	if IsNotLoaded(err) {
		result, err = dr.QueryDNSRecordToTag().All(ctx)
	}
	return result, err
}

func (d *Disk) DiskToTag(ctx context.Context) ([]*Tag, error) {
	result, err := d.Edges.DiskToTagOrErr()
	if IsNotLoaded(err) {
		result, err = d.QueryDiskToTag().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToTag(ctx context.Context) ([]*Tag, error) {
	result, err := e.Edges.EnvironmentToTagOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToTag().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToUser(ctx context.Context) ([]*User, error) {
	result, err := e.Edges.EnvironmentToUserOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToUser().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToHost(ctx context.Context) ([]*Host, error) {
	result, err := e.Edges.EnvironmentToHostOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToHost().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToCompetition(ctx context.Context) ([]*Competition, error) {
	result, err := e.Edges.EnvironmentToCompetitionOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToCompetition().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToBuild(ctx context.Context) ([]*Build, error) {
	result, err := e.Edges.EnvironmentToBuildOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToBuild().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToIncludedNetwork(ctx context.Context) ([]*IncludedNetwork, error) {
	result, err := e.Edges.EnvironmentToIncludedNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToIncludedNetwork().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToNetwork(ctx context.Context) ([]*Network, error) {
	result, err := e.Edges.EnvironmentToNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToNetwork().All(ctx)
	}
	return result, err
}

func (e *Environment) EnvironmentToTeam(ctx context.Context) ([]*Team, error) {
	result, err := e.Edges.EnvironmentToTeamOrErr()
	if IsNotLoaded(err) {
		result, err = e.QueryEnvironmentToTeam().All(ctx)
	}
	return result, err
}

func (fd *FileDelete) FileDeleteToTag(ctx context.Context) ([]*Tag, error) {
	result, err := fd.Edges.FileDeleteToTagOrErr()
	if IsNotLoaded(err) {
		result, err = fd.QueryFileDeleteToTag().All(ctx)
	}
	return result, err
}

func (fd *FileDownload) FileDownloadToTag(ctx context.Context) ([]*Tag, error) {
	result, err := fd.Edges.FileDownloadToTagOrErr()
	if IsNotLoaded(err) {
		result, err = fd.QueryFileDownloadToTag().All(ctx)
	}
	return result, err
}

func (fe *FileExtract) FileExtractToTag(ctx context.Context) ([]*Tag, error) {
	result, err := fe.Edges.FileExtractToTagOrErr()
	if IsNotLoaded(err) {
		result, err = fe.QueryFileExtractToTag().All(ctx)
	}
	return result, err
}

func (f *Finding) FindingToUser(ctx context.Context) ([]*User, error) {
	result, err := f.Edges.FindingToUserOrErr()
	if IsNotLoaded(err) {
		result, err = f.QueryFindingToUser().All(ctx)
	}
	return result, err
}

func (f *Finding) FindingToTag(ctx context.Context) ([]*Tag, error) {
	result, err := f.Edges.FindingToTagOrErr()
	if IsNotLoaded(err) {
		result, err = f.QueryFindingToTag().All(ctx)
	}
	return result, err
}

func (f *Finding) FindingToHost(ctx context.Context) ([]*Host, error) {
	result, err := f.Edges.FindingToHostOrErr()
	if IsNotLoaded(err) {
		result, err = f.QueryFindingToHost().All(ctx)
	}
	return result, err
}

func (f *Finding) FindingToScript(ctx context.Context) ([]*Script, error) {
	result, err := f.Edges.FindingToScriptOrErr()
	if IsNotLoaded(err) {
		result, err = f.QueryFindingToScript().All(ctx)
	}
	return result, err
}

func (h *Host) HostToDisk(ctx context.Context) ([]*Disk, error) {
	result, err := h.Edges.HostToDiskOrErr()
	if IsNotLoaded(err) {
		result, err = h.QueryHostToDisk().All(ctx)
	}
	return result, err
}

func (h *Host) HostToUser(ctx context.Context) ([]*User, error) {
	result, err := h.Edges.HostToUserOrErr()
	if IsNotLoaded(err) {
		result, err = h.QueryHostToUser().All(ctx)
	}
	return result, err
}

func (h *Host) HostToTag(ctx context.Context) ([]*Tag, error) {
	result, err := h.Edges.HostToTagOrErr()
	if IsNotLoaded(err) {
		result, err = h.QueryHostToTag().All(ctx)
	}
	return result, err
}

func (h *Host) HostToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := h.Edges.HostToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = h.QueryHostToEnvironment().All(ctx)
	}
	return result, err
}

func (in *IncludedNetwork) IncludedNetworkToTag(ctx context.Context) ([]*Tag, error) {
	result, err := in.Edges.IncludedNetworkToTagOrErr()
	if IsNotLoaded(err) {
		result, err = in.QueryIncludedNetworkToTag().All(ctx)
	}
	return result, err
}

func (in *IncludedNetwork) IncludedNetworkToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := in.Edges.IncludedNetworkToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = in.QueryIncludedNetworkToEnvironment().All(ctx)
	}
	return result, err
}

func (n *Network) NetworkToTag(ctx context.Context) ([]*Tag, error) {
	result, err := n.Edges.NetworkToTagOrErr()
	if IsNotLoaded(err) {
		result, err = n.QueryNetworkToTag().All(ctx)
	}
	return result, err
}

func (n *Network) NetworkToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := n.Edges.NetworkToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = n.QueryNetworkToEnvironment().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToTag(ctx context.Context) ([]*Tag, error) {
	result, err := ph.Edges.ProvisionedHostToTagOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToTag().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToStatus(ctx context.Context) ([]*Status, error) {
	result, err := ph.Edges.ProvisionedHostToStatusOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToStatus().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToProvisionedNetwork(ctx context.Context) ([]*ProvisionedNetwork, error) {
	result, err := ph.Edges.ProvisionedHostToProvisionedNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToProvisionedNetwork().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToHost(ctx context.Context) ([]*Host, error) {
	result, err := ph.Edges.ProvisionedHostToHostOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToHost().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToProvisioningStep(ctx context.Context) ([]*ProvisioningStep, error) {
	result, err := ph.Edges.ProvisionedHostToProvisioningStepOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToProvisioningStep().All(ctx)
	}
	return result, err
}

func (ph *ProvisionedHost) ProvisionedHostToAgentStatus(ctx context.Context) ([]*AgentStatus, error) {
	result, err := ph.Edges.ProvisionedHostToAgentStatusOrErr()
	if IsNotLoaded(err) {
		result, err = ph.QueryProvisionedHostToAgentStatus().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToTag(ctx context.Context) ([]*Tag, error) {
	result, err := pn.Edges.ProvisionedNetworkToTagOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToTag().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToStatus(ctx context.Context) ([]*Status, error) {
	result, err := pn.Edges.ProvisionedNetworkToStatusOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToStatus().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToNetwork(ctx context.Context) ([]*Network, error) {
	result, err := pn.Edges.ProvisionedNetworkToNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToNetwork().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToBuild(ctx context.Context) ([]*Build, error) {
	result, err := pn.Edges.ProvisionedNetworkToBuildOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToBuild().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToTeam(ctx context.Context) ([]*Team, error) {
	result, err := pn.Edges.ProvisionedNetworkToTeamOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToTeam().All(ctx)
	}
	return result, err
}

func (pn *ProvisionedNetwork) ProvisionedNetworkToProvisionedHost(ctx context.Context) ([]*ProvisionedHost, error) {
	result, err := pn.Edges.ProvisionedNetworkToProvisionedHostOrErr()
	if IsNotLoaded(err) {
		result, err = pn.QueryProvisionedNetworkToProvisionedHost().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToTag(ctx context.Context) ([]*Tag, error) {
	result, err := ps.Edges.ProvisioningStepToTagOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToTag().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToStatus(ctx context.Context) ([]*Status, error) {
	result, err := ps.Edges.ProvisioningStepToStatusOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToStatus().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToProvisionedHost(ctx context.Context) ([]*ProvisionedHost, error) {
	result, err := ps.Edges.ProvisioningStepToProvisionedHostOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToProvisionedHost().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToScript(ctx context.Context) ([]*Script, error) {
	result, err := ps.Edges.ProvisioningStepToScriptOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToScript().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToCommand(ctx context.Context) ([]*Command, error) {
	result, err := ps.Edges.ProvisioningStepToCommandOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToCommand().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToDNSRecord(ctx context.Context) ([]*DNSRecord, error) {
	result, err := ps.Edges.ProvisioningStepToDNSRecordOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToDNSRecord().All(ctx)
	}
	return result, err
}

func (ps *ProvisioningStep) ProvisioningStepToRemoteFile(ctx context.Context) ([]*RemoteFile, error) {
	result, err := ps.Edges.ProvisioningStepToRemoteFileOrErr()
	if IsNotLoaded(err) {
		result, err = ps.QueryProvisioningStepToRemoteFile().All(ctx)
	}
	return result, err
}

func (rf *RemoteFile) RemoteFileToTag(ctx context.Context) ([]*Tag, error) {
	result, err := rf.Edges.RemoteFileToTagOrErr()
	if IsNotLoaded(err) {
		result, err = rf.QueryRemoteFileToTag().All(ctx)
	}
	return result, err
}

func (s *Script) ScriptToTag(ctx context.Context) ([]*Tag, error) {
	result, err := s.Edges.ScriptToTagOrErr()
	if IsNotLoaded(err) {
		result, err = s.QueryScriptToTag().All(ctx)
	}
	return result, err
}

func (s *Script) ScriptToUser(ctx context.Context) ([]*User, error) {
	result, err := s.Edges.ScriptToUserOrErr()
	if IsNotLoaded(err) {
		result, err = s.QueryScriptToUser().All(ctx)
	}
	return result, err
}

func (s *Script) ScriptToFinding(ctx context.Context) ([]*Finding, error) {
	result, err := s.Edges.ScriptToFindingOrErr()
	if IsNotLoaded(err) {
		result, err = s.QueryScriptToFinding().All(ctx)
	}
	return result, err
}

func (s *Status) StatusToTag(ctx context.Context) ([]*Tag, error) {
	result, err := s.Edges.StatusToTagOrErr()
	if IsNotLoaded(err) {
		result, err = s.QueryStatusToTag().All(ctx)
	}
	return result, err
}

func (t *Team) TeamToUser(ctx context.Context) ([]*User, error) {
	result, err := t.Edges.TeamToUserOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryTeamToUser().All(ctx)
	}
	return result, err
}

func (t *Team) TeamToBuild(ctx context.Context) ([]*Build, error) {
	result, err := t.Edges.TeamToBuildOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryTeamToBuild().All(ctx)
	}
	return result, err
}

func (t *Team) TeamToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := t.Edges.TeamToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryTeamToEnvironment().All(ctx)
	}
	return result, err
}

func (t *Team) TeamToTag(ctx context.Context) ([]*Tag, error) {
	result, err := t.Edges.TeamToTagOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryTeamToTag().All(ctx)
	}
	return result, err
}

func (t *Team) TeamToProvisionedNetwork(ctx context.Context) ([]*ProvisionedNetwork, error) {
	result, err := t.Edges.TeamToProvisionedNetworkOrErr()
	if IsNotLoaded(err) {
		result, err = t.QueryTeamToProvisionedNetwork().All(ctx)
	}
	return result, err
}

func (u *User) UserToTag(ctx context.Context) ([]*Tag, error) {
	result, err := u.Edges.UserToTagOrErr()
	if IsNotLoaded(err) {
		result, err = u.QueryUserToTag().All(ctx)
	}
	return result, err
}

func (u *User) UserToEnvironment(ctx context.Context) ([]*Environment, error) {
	result, err := u.Edges.UserToEnvironmentOrErr()
	if IsNotLoaded(err) {
		result, err = u.QueryUserToEnvironment().All(ctx)
	}
	return result, err
}
