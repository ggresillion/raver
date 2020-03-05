import {Component} from '@angular/core';
import {MatDialogRef} from '@angular/material/dialog';
import {environment} from '../../../../environments/environment';
import * as querystring from 'querystring';

@Component({
  selector: 'app-add-bot-to-guild-dialog',
  templateUrl: './add-bot-to-guild-dialog.component.html',
  styleUrls: ['./add-bot-to-guild-dialog.component.scss']
})
export class AddBotToGuildDialogComponent {

  constructor(public readonly dialogRef: MatDialogRef<AddBotToGuildDialogComponent>) {
  }

  public navigateToLink(): void {
    window.open(`${environment.discord.api}?${querystring.encode({
      client_id: environment.discord.clientId,
      permissions: environment.discord.permissions,
      redirect_uri: location.href,
      scope: environment.discord.scope,
    })}`);
    this.dialogRef.close();
  }
}
