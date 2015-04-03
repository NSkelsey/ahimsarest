'use strict';

angular.module('ombWebAppFactory', []).factory('ahimsaRestService', function($http) {
  // the endpoint must be from the same origin, otherwise this doesn't work!
  return {
    'getAllBoards': function() {
      return $http.get('/api/boards')
    },
    'getBoard': function(urlname) {
    // Use the encode function to prevent user submitted data from doing funky stuff with the url.
      return $http.get('/api/board/' + encodeURIComponent(urlname))
    },
    'getNilBoard': function() {
      return $http.get('/api/nilboard')
    },
 }
});


