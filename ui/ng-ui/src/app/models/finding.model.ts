import { Host } from '@angular/core';
import { Tag, User } from './common.model';

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
  // maintainer: User;
  // tags: Tag[];
  // Host: Host;
  findingToUser: User;
  findingToTag: Tag[];
  findingToHost: Host;
}
