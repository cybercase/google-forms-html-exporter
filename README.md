# Google Forms Exporter
## Convert any Google Form to an HTML form

_Live @_ https://stefano.brilli.me/google-forms-html-exporter/



## Developers area

The project has 2 parts:

- backend
- frontend

### Building the backend
cd `cmd/formdress` && go build

### Building the frontend
You'll need version 8 of [node](https://nodejs.org/), [bower](https://bower.io/) and [npm](https://www.npmjs.com/).

run `npm install`, `bower install` then `./node_modules/.bin/gulp` to build the frontend

### Run on localhost

- Change the server address in `app/scripts/config.js` to `http://localhost:8000`
- Build the backend, build the frontend, then run `./cmd/formdress/formdress -d ./docs`.
- Point your browser to `http://localhost:8000`

### Using as tool
You can also use the `./cmd/formdress/formdress` command as a local tool to export Google Forms as json objects.
Just type: `./cmd/formdress/formdress -f [YOUR_GOOGLE_FORM_URL]`
