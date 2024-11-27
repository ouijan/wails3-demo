import { Component, signal } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
import { GreetService } from '@wails/app';
import { Events } from '@wailsio/runtime';
import { get } from 'lodash';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, ReactiveFormsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'ng-app';
  greeting = signal('');
  timeDisplay = signal('');
  nameControl = new FormControl('');

  constructor() {
    Events.On('time', (time: unknown) => {
      const data = get(time, 'data', '') as string;
      this.timeDisplay.set(data);
    });

    setInterval(() => {
      const date = new Date().toISOString();
      console.log(`SyncCheck: ${date}`);
      GreetService.SyncCheck(date);
    }, 1000);
  }

  async submit(): Promise<void> {
    const name = this.nameControl.value;
    if (!name) {
      return;
    }
    const response = await GreetService.Greet(name);
    this.greeting.set(response);
  }
}
