import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'dateBetween',
  pure: true
})
export class DateBetweenPipe implements PipeTransform {
  transform(firstDate: any, secondDate: any): any {
    if (firstDate && secondDate) {
      let seconds = Math.floor((+new Date(secondDate) - +new Date(firstDate)) / 1000);
      if (seconds < 0) {
        seconds = Math.floor((+new Date(firstDate) - +new Date(secondDate)) / 1000);
      }
      // if (seconds < 29)
      //   // less than 30 seconds ago will show as 'Just now'
      //   return 'Just now';
      const intervals = {
        year: 31536000,
        month: 2592000,
        week: 604800,
        day: 86400,
        hour: 3600,
        minute: 60,
        second: 1
      };
      let counter;
      for (const i in intervals) {
        counter = Math.floor(seconds / intervals[i]);
        if (counter > 0)
          if (counter === 1) {
            return counter + ' ' + i; // singular (1 day ago)
          } else {
            return counter + ' ' + i + 's'; // plural (2 days ago)
          }
      }
    }
    return firstDate;
  }
}
