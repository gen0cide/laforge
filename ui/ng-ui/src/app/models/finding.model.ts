import { Host } from '@angular/core';
import { Tag, User } from './common.model';

enum FindingSeverity {
  ZeroSeverity,
  LowSeverity,
  MediumSeverity,
  HighSeverity,
  CriticalSeverity,
  NullSeverity
}

enum FindingDifficulty {
  ZeroDifficulty,
  NoviceDifficulty,
  AdvancedDifficulty,
  ExpertDifficulty,
  NullDifficulty
}

interface Finding {
  name: string;
  description: string;
  severity: FindingSeverity;
  difficulty: FindingDifficulty;
  maintainer: User;
  tags: Tag[];
  Host: Host;
}

export { Finding };
