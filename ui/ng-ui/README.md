# Laforge UI

This is the UI frontend to interact with the new Laforge Server.

## Development server

Duplicate the `src/environments/environment.ts` to the file `src/environments/environment.dev.ts` and modify the configuration as needed.

Run `ng serve --configuration=dev` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

```
Components - src/app/components
Pages      - src/app/pages
Services   - src/app/services
Pipes      - src/app/pipes
```

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory.z

### Production

When building the UI for production, use:

```
ng build --prod
```

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
