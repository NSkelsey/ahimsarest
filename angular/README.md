Angular-web
===========

After coming back to this project after a few months, I have written some notes to 
remind myself and others as to how this little web-app works.

### Running 

This is a go web server so as usual it can be run with:
```bash
> go run app.go -statpath="./" -host="0.0.0.0:1055"
```

This creates a webserver that serves files at http://0.0.0.0:1055/ and an ombuds json api at http://0.0.0.0:1055/api/.

### Understanding

Under the hood this program is two things. It is an angular frontend and a json api. 
The routes to each are documented above.
An angular app is a single page web application that consumes an api to display and manipulate data.
Angular has a built in templating syntax along with abstractions that make it easy to consume an api like the one we created.

The json objects that are reachable through the api are documented [here](https://godoc.org/github.com/soapboxsys/ombudslib/ombjson).
The api itself is really only these [few lines](https://github.com/NSkelsey/ahimsarest/blob/master/jsonapi.go#L270-L286) 
(Note that the prefix in the angular example is `/api/`).


`index.html` is really the starting point of the application. 
It does several things (go ahead and open it now).
Notice the ng-app="ahimsaApp" at the top. 
ng-[name] are attribute tags that tell angular's javascript what to attach itself to.

Notice that portions of the page are broken into controllers. 
The top bar has a `headCtrl` and their is an `<ngview>` that sits in the middles of the page.

Finally at the bottom of the page the javascript gets loaded and executed. 
Since we have defined additional custom views along with a main.js, those come last.

### Gotchas

Since it is a single page app, if the server goes down the api will stop working, 
but the links will still work as long as you stay on that one page. 


### Todo
The first thing you should do is implement a status indicator that checks the health of the api.
