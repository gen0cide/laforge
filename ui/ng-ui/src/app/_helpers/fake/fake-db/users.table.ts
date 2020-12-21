export class UsersTable {
  public static users: any = [
    {
      id: 1,
      username: 'lucas',
      password: 'admin',
      email: 'lucas@cptc.test',
      accessToken: 'access-token-8f3ae836da744329a6f93bf20594b5cc',
      refreshToken: 'access-token-f8c137a2c98743f48b643e71161d90aa',
      roles: [1], // Administrator
      pic: './assets/media/users/300_25.jpg',
      fullname: 'Lucas',
      occupation: 'Developer',
      companyName: 'CPTC',
      phone: '1234567890',
      address: {
        addressLine: '',
        city: '',
        state: '',
        postCode: '12345',
      },
      socialNetworks: {
        linkedIn: 'https://linkedin.com/admin',
        facebook: 'https://facebook.com/admin',
        twitter: 'https://twitter.com/admin',
        instagram: 'https://instagram.com/admin',
      },
    },
    {
      id: 2,
      username: 'user',
      password: 'demo',
      email: 'user@demo.com',
      accessToken: 'access-token-6829bba69dd3421d8762-991e9e806dbf',
      refreshToken: 'access-token-f8e4c61a318e4d618b6c199ef96b9e55',
      roles: [2], // Manager
      pic: './assets/media/users/100_2.jpg',
      fullname: 'Megan',
      occupation: 'Deputy Head of Keenthemes in New York office',
      companyName: 'Keenthemes',
      phone: '456669067891',
      address: {
        addressLine: '3487  Ingram Road',
        city: 'Greensboro',
        state: 'North Carolina',
        postCode: '27409',
      },
      socialNetworks: {
        linkedIn: 'https://linkedin.com/user',
        facebook: 'https://facebook.com/user',
        twitter: 'https://twitter.com/user',
        instagram: 'https://instagram.com/user',
      },
    },
    {
      id: 3,
      username: 'guest',
      password: 'demo',
      email: 'guest@demo.com',
      accessToken: 'access-token-d2dff7b82f784de584b60964abbe45b9',
      refreshToken: 'access-token-c999ccfe74aa40d0aa1a64c5e620c1a5',
      roles: [3], // Guest
      pic: './assets/media/users/default.jpg',
      fullname: 'Ginobili Maccari',
      occupation: 'CFO',
      companyName: 'Keenthemes',
      phone: '456669067892',
      address: {
        addressLine: '1467  Griffin Street',
        city: 'Phoenix',
        state: 'Arizona',
        postCode: '85012',
      },
      socialNetworks: {
        linkedIn: 'https://linkedin.com/guest',
        facebook: 'https://facebook.com/guest',
        twitter: 'https://twitter.com/guest',
        instagram: 'https://instagram.com/guest',
      },
    },
  ];

  public static tokens: any = [
    {
      id: 1,
      accessToken: 'access-token-' + Math.random(),
      refreshToken: 'access-token-' + Math.random(),
    },
  ];
}
