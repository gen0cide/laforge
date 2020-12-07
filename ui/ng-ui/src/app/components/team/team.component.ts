import {
  ChangeDetectorRef,
  Component,
  ElementRef,
  Input,
  OnInit,
  Renderer2,
  ViewChild
} from '@angular/core';
import { Team } from 'src/app/models/common.model';

@Component({
  selector: 'app-team',
  templateUrl: './team.component.html',
  styleUrls: ['./team.component.scss']
})
export class TeamComponent implements OnInit {
  @Input() team: Team;

  constructor() {}

  ngOnInit(): void {}

  getStatus(): string {
    return 'ProvStatusInProgress';
  }
}
