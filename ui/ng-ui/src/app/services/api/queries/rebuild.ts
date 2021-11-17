import { ID } from '@models/common.model';
import { gql } from 'apollo-angular';

export interface RebuildPlansVars {
  rootPlans: ID[];
}

export interface RebuildPlansData {
  rebuild: boolean;
}

export const RebuildPlansMutation = gql`
  mutation($rootPlans: [String]!) {
    rebuild(rootPlans: $rootPlans)
  }
`;
