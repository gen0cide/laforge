import { tagMap, User } from './common.model';
import { Environment } from './environment.model';
import { Script } from './script.model';

export enum FindingSeverity {
  ZeroSeverity,
  LowSeverity,
  MediumSeverity,
  HighSeverity,
  CriticalSeverity,
  NullSeverity
}

export enum FindingDifficulty {
  ZeroDifficulty,
  NoviceDifficulty,
  AdvancedDifficulty,
  ExpertDifficulty,
  NullDifficulty
}

export interface Finding {
  name: string;
  description: string;
  severity: FindingSeverity;
  difficulty: FindingDifficulty;
  tags: tagMap[];
  FindingToUser: User[];
  FindingToScript: Script;
  FindingToEnvironment: Environment;
}
