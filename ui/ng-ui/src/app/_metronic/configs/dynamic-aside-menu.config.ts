export const DynamicAsideMenuConfig = {
  items: [
    {
      title: 'Dashboard',
      root: true,
      icon: 'flaticon2-architecture-and-city',
      svg: './assets/media/svg/icons/Home/Home.svg',
      page: '/dashboard',
      translate: 'MENU.DASHBOARD',
      bullet: 'dot',
    },
    { section: 'Develop' },
    {
      title: 'Plan',
      root: true,
      bullet: 'dot',
      icon: 'flaticon2-list-2',
      svg: './assets/media/svg/icons/Design/Layers.svg',
      page: '/develop/plan',
    },
    {
      title: 'Build',
      root: true,
      bullet: 'dot',
      icon: 'flaticon2-list-2',
      svg: './assets/media/svg/icons/Tools/Hummer2.svg',
      page: '/develop/plan',
    },
    { section: 'Manage' },
    {
      title: 'Environments',
      root: true,
      bullet: 'dot',
      icon: 'flaticon2-list-2',
      svg: './assets/media/svg/icons/Devices/Server.svg',
      page: '/develop/plan',
    }
  ]
};
