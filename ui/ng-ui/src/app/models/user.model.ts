import { ID } from './common.model';

export interface Role {
  id: ID;
  name: string;
  // permissions field if we want roles to act as "groups"
}

export interface User {
  id: ID;
  username: string;
  email: string;
  accessToken: string; // consider deprecating these fields in favor of cookies set by auth server
  refreshToken: string; // ^^^^
  roles: Role[];
  profilePicture: string; // CDN URL to user profile picture
  firstName: string;
  lastName: string;
  occupation: string; // probably not needed? use roles to further explain
  companyName: string; // possible convert to a region object?
  phone: string;
  // address: Address; do we really want to store user addresses?
  // socialNetworks field if we really really want it (store as it's own object or limit to certain ones?)
}
