# Google Forms Exporter
## Convert any Google Form to an HTML form

_Live @_ https://stefano.brilli.me/google-forms-html-exporter/



## Developers area

The project has 2 parts:

- backend
- frontend


### Building the backend
you'll need [go](https://golang.org/) and [gb](https://github.com/constabulary/gb)

run `gb vendor restore` to install all the dependencies, then `gb build` to build the backend (output to `./bin` folder).

### Building the frontend
you'll need [node](https://nodejs.org/), [bower](https://bower.io/) and [npm](https://www.npmjs.com/).

run `npm install`, `bower install` then `./node_modules/.bin/gulp` to build the frontend

### Run on localhost

- Change the server address in `app/scripts/config.js` to `http://localhost:8000`
- Build the backend, build the frontend, then run `./bin/formdress -d ./docs`.
- Point your browser to `http://localhost:8000`

### Using as tool
You can also use the `./bin/formdress` command as a local tool to export Google Forms as json objects.
Just type: `./bin/formdress -f [YOUR_GOOGLE_FORM_URL]`
