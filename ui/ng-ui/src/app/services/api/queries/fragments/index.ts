import { gql } from 'apollo-angular';

const StatusFields = gql`
  fragment StatusFields on Status {
    state
    started_at
    ended_at
    failed
    completed
    error
  }
`;

const PlanFields = gql`
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
