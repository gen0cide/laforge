
/* RENAME ME TO "users.table.prod.ts" */

export class UsersTable {
  public static users: any = [
    {
      id: 1,
      username: 'user',
      password: 'Password123',
      email: 'user@example.com',
      accessToken: 'access-token-a3ad6f09536498c72bc85fb443f32ea9',
      refreshToken: 'access-token-314819a98c0d66efa7a117738324cb4f',
      roles: [1], // Administrator
      pic: './assets/media/users/300_21.jpg',
      fullname: 'User',
      occupation: 'Developer',
      companyName: 'CPTC',
      phone: '1234567890',
      address: {
        addressLine: '123 Main St.',
        city: 'Fakecity',
        state: 'AA',
        postCode: '12345',
      },
      socialNetworks: {
        linkedIn: 'https://linkedin.com/user',
        facebook: 'https://facebook.com/user',
        twitter: 'https://twitter.com/user',
        instagram: 'https://instagram.com/admuserin',
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
