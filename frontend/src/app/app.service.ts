import { Injectable } from '@angular/core';
import { GreetService } from '../../bindings/github.com/ouijan/wails3-demo/backend/app';
import { Events } from '@wailsio/runtime';

@Injectable()
export class WailsBindings {
  Events = Events;
  GreetService = GreetService;
}
