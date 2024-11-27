import { Component, inject, signal } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
import { get } from 'lodash';
import { WailsBindings } from './app.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, ReactiveFormsModule],
  providers: [WailsBindings],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  bindings = inject(WailsBindings);
  title = 'ng-app';
  greeting = signal('');
  timeDisplay = signal('');
  nameControl = new FormControl('');

  constructor() {
    this.bindings.Events.On('time', (time: unknown) => {
      const data = get(time, 'data', '') as string;
      this.timeDisplay.set(data);
    });

    setInterval(() => {
      const date = new Date().toISOString();
      console.log(`SyncCheck: ${date}`);
      this.bindings.GreetService.SyncCheck(date);
    }, 1000);
  }

  async submit(): Promise<void> {
    const name = this.nameControl.value;
    if (!name) {
      return;
    }
    const response = await this.bindings.GreetService.Greet(name);
    this.greeting.set(response);
  }
}
