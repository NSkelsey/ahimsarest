'use strict';

function formatDayStr(date) {
    var m = String(date.getUTCMonth() + 1);
    var day = String(date.getDate());
    var year = date.getUTCFullYear();
    var monthStr = m.length > 1 ? m : "0" + m;
    var dayStr = day.length > 1 ? day : "0" + day;

    return "{0}-{1}-{2}".format(dayStr, monthStr, year);
}

angular.module('ahimsaApp', [
    'ngRoute', 
    'ahimsaApp.blockView',
    'ahimsaApp.blocksView',
    ])
    .config(['$routeProvider', function($routeProvider) {
              $routeProvider.when('/', {
                  controller: 'allboards', 
                  templateUrl: 'boards.html'
              })
              .when('/boards', {
                  controller: 'allboards',
                  templateUrl: 'boards.html'
              })
              .when('/board/:board*', {
                  controller: 'board', 
                  templateUrl: 'board.html'
              })
              .when('/nilboard', {
                  controller: 'nilboard', 
                  templateUrl: 'board.html'
              })
              .when('/bulletin/:txid', {
                  controller: 'bulletin', 
                  templateUrl: 'bltn.html'
              })
              .when('/author/:author', {
                  controller: 'author',
                  templateUrl: 'author.html'
              })
              .when('/authors', {
                  controller: 'authors',
                  templateUrl: 'authors.html'
              })
              .when('/blacklist', {
                  controller: 'blacklist',
                  templateUrl: 'blacklist.html'
              })
              .when('/blacklist/:txid', {
                  controller: 'blacklist',
                  templateUrl: 'blacklist.html'
              })
              $routeProvider.otherwise({
                  redirectTo: '/'
              });
    }])
    .factory('ahimsaRestService', function($http) {
        // the endpoint must be from the same origin, otherwise this doesn't work!
        return {
            'getAllBoards': function() {
                return $http.get('/api/boards')
            },
            'getBoard': function(urlname) {
                return $http.get('/api/board/' + encodeURIComponent(urlname))
            },
            'getNilBoard': function() {
                return $http.get('/api/nilboard')
            },
            'getBulletin': function(txid) {
                return $http.get('/api/bulletin/' + txid)
            },
            'getAuthor': function(author) {
                return $http.get('/api/author/' + author)
            },
            'getAllAuthors': function() {
                return $http.get('/api/authors')
            },
            'getStatus': function() {
                return $http.get('/api/status')
            },
            'getBlacklist': function() {
                return $http.get('/api/blacklist')
            },
            'getBlocksDay': function(day) {
                return $http.get('/api/blocks/' + day)
            },
            'getBlock': function(hash) {
                return $http.get('/api/block/' + hash)
            },
        }
    })
    .controller('headCtrl', [
      // To minify controllers annotate variables
      '$scope', 
      '$location', 
      'ahimsaRestService', 
      function($scope, $location, ahimsaRestService){
        $scope.isActive = function(route) {
            // Determine if location.path starts with 'route'
            return $location.path().indexOf(route, 0) === 0
        }

        // Set today's date
        $scope.today = ((new Date()).getTime() / 1000)

        ahimsaRestService.getStatus().then(function(result) {
            $scope.status = result.data
        });
    }])
    .controller('allboards', function($scope, ahimsaRestService) {

        ahimsaRestService.getAllBoards().then(function(result) {
            $scope.boards = angular.forEach(result.data, initBoard);
        });

    })
    .controller('board', function($scope, $routeParams, ahimsaRestService) {

        ahimsaRestService.getBoard($routeParams.board).then(function(result) {
            $scope.board = initBoard(result.data);
        });
    })
    .controller('nilboard', function($scope, ahimsaRestService) {
        ahimsaRestService.getNilBoard().then(function(result) {
            $scope.board = result.data;
            $scope.board.summary.name = "The nil board."
        });
    })
    .controller('bulletin', function($scope, $routeParams, ahimsaRestService) {
        ahimsaRestService.getBulletin($routeParams.txid).then(function(result) {
            $scope.bltn = result.data;
        });
    })
    .controller('author', function($scope, $routeParams, ahimsaRestService) {
        ahimsaRestService.getAuthor($routeParams.author).then(function(result) {
            $scope.authSum = result.data
        });
    })
    .controller('authors', function($scope, ahimsaRestService) {
        ahimsaRestService.getAllAuthors().then(function(result) {
            $scope.authors = result.data
        });
    })
    .controller('blacklist', function($scope, $routeParams, ahimsaRestService) {
        if ($routeParams.txid !== undefined) {
            $scope.active = $routeParams.txid
        }

        ahimsaRestService.getBlacklist().then(function(result) {
            $scope.blacklist = result.data 
        });
    })
    .controller('status', function($scope, ahimsaRestService) {
        ahimsaRestService.getStatus().then(function(result) {
            $scope.status = result.data
        });
    })

    // Directives and filters
    .directive('bulletin', function() {
        return {
            templateUrl: 'bltn.html',
            restrict: 'C',
        }
    })
    .directive('author', function() {
        return {
            templateUrl: 'author-line.html',
            restrict: 'C',
        }
    })
    .directive('block', function() {
        return {
            templateUrl: 'block-line.html',
            restrict: 'C',
        }
    })

    .filter('epochdate', function() {
            return function(utcsecs) {
                return new Date(utcsecs*1000)
            }
    })

    .filter('precisedate', function() {
        var options = { 
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: 'numeric',  
            minute: 'numeric',
            timeZone: 'UTC',
        };
        var formater = new Intl.DateTimeFormat('en-GB', options);
        return function(utcsecs) {
            var d = new Date(utcsecs*1000); 
            return formater.format(d);
        }
    })
    .filter('nicedate', function() {
        return function(date) {
            var options = { year: 'numeric', month: 'short', day: 'numeric' }
            return date.toLocaleDateString('en-US', options)
        }
    })
    .filter('formatdaystr', function() {
        return function(date) {
            var m = String(date.getUTCMonth() + 1);
            var day = String(date.getDate());
            var year = date.getUTCFullYear();
            var monthStr = m.length > 1 ? m : "0" + m;
            var dayStr = day.length > 1 ? day : "0" + day;

            return "{0}-{1}-{2}".format(dayStr, monthStr, year);
        }
    });



// Takes in a board and adds an urlname
function initBoard(board) {
    board.urlName = encodeURIComponent(board.name)
    return board
}

// Takes in a daystr and spits out a date object
function parseDayStr(str) {
    var re = /([0-9]{1,2})-([0-9]{1,2})-([0-9]{4})/
    var m = str.match(re)
    if (m === null) {
        return "I Broke"
    }
    else {
        return new Date(m[3], (m[2] - 1), m[1], 0, 0, 0)
    }
}

