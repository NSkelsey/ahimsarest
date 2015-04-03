var ctrls = angular.module('ombWebAppControllers', []);

ctrls.controller('board', function($scope, $routeParams, ahimsaRestService) {
      ahimsaRestService.getBoard($routeParams.board).then(function(result) {
          initBoardSum(result.data.summary);
          angular.forEach(result.data.bltns, initBulletin);
          $scope.board = result.data;
  });
})
  
ctrls.controller('welcome', function($scope) {
    
});

ctrls.controller('browseCtrl', function($scope, $location, $routeParams, ahimsaRestService) {
    ahimsaRestService.getAllBoards().then(function(result) {
        $scope.boards = angular.forEach(result.data, initBoardSum);
        $scope.openBoard = function(name) {
           $location.path("/board/" + name) 
        }
    });
})


// Takes in a board and adds processed fields
function initBoardSum(board) {
    board.urlName = encodeURIComponent(board.name);
    board.lastActiveDate = simpleDate(board.lastActive);
    board.createdAtDateTime = simpleDateTime(board.createdAt);
    return board;
}

function initBulletin(bltn) {
    bltn.timestampDateTime = simpleDateTime(bltn.timestamp);
}

function simpleDate(utcsecs) {
    var d = new Date(utcsecs*1000); 
    var options = { year: '2-digit', month: '2-digit', day: '2-digit' }
    return d.toLocaleDateString('en-US', options)
}

var options = { 
    year: '2-digit',
    month: '2-digit',
    day: '2-digit',
    hour: 'numeric',  
    minute: 'numeric',
    timeZone: 'UTC',
};

var formater = new Intl.DateTimeFormat('en-US', options);

function simpleDateTime(utcsecs) {
    var d = new Date(utcsecs*1000); 
    return formater.format(d);
}
