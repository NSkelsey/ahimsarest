'use strict';

var ombWebApp = angular.module("ombWebApp", [
    'ombWebAppControllers',
    'ombWebAppFilters',
    'ombWebAppFactory',
    'ngRoute',
    ])
    .config(["$routeProvider", function($routeProvider) {
      $routeProvider.when('/', {
        controller: 'welcome',
        templateUrl: 'welcome.html'
      })
      .when('/board/:board*', {
        controller: 'board',
        templateUrl: 'board.html'
      })
    }])
    // Directives
    .directive('browseBoard', function() {
        return {
            templateUrl: 'browse-board.html',
            // the restrict tag means that only html elements that have the class "bulletin"
            // end up with the bltn.html template rendered below them.
            restrict: 'C',
        }
    })
    .directive('pinnedBulletin', function() {
        return {
            templateUrl: 'pinned-bulletin.html', 
            //restrict: 'C',
        }
    })
