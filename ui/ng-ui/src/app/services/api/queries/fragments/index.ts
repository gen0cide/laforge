import { gql } from 'apollo-angular';

export const StatusFields = gql`
  fragment StatusFields on Status {
    state
    started_at
    ended_at
    failed
    completed
    error
  }
`;

export const PlanFields = gql`
  fragment PlanFields on Plan {
    id
    step_number
    type
    PlanToStatus {
      ...StatusFields
    }
  }

  ${StatusFields}
`;
