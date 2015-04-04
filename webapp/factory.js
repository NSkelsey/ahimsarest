'use strict';

angular.module('ombWebAppFactory', []).factory('ahimsaRestService', function($http, $interval) {
  
  var serverInfo = { 'blkHeight': 0 };

  function updateService() { 
      $http.get('/api/status').then(function(result) {
          serverInfo = result.data;
          debugger;
      });
  }

  updateService();
  $interval(updateService, 15*1000, 0, true);
 
  
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
    'getBlockHeight': function() {
        return serverInfo.blkHeight;
    },
 }
})


