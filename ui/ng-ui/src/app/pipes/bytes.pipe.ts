import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'fromBytes' })
export class FromBytesPipe implements PipeTransform {
  transform(value: number, precision?: number): number | string {
    if (isNaN(value) || !isFinite(value)) return '-';
    if (typeof precision === 'undefined') precision = 1;
    const units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'];
    const number = Math.floor(Math.log(value) / Math.log(1024));
    return (value / Math.pow(1024, Math.floor(number))).toFixed(precision) + ' ' + units[number];
  }
}

// app.filter('bytes', function() {
// 	return function(bytes, precision) {
// 		if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
// 		if (typeof precision === 'undefined') precision = 1;
// 		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
// 			number = Math.floor(Math.log(bytes) / Math.log(1024));
// 		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
// 	}
// });
