var ctrls = angular.module('ombWebAppControllers', []);

ctrls.controller('board', function($scope, $routeParams, ahimsaRestService) {
      ahimsaRestService.getBoard($routeParams.board).then(function(result) {
          var board = result.data;
          $scope.board = board;
  });
})

ctrls.controller('nilboard', function($scope, ahimsaRestService) {
      ahimsaRestService.getNilBoard().then(function(result) {
          $scope.board = result.data;
  });
})
  
ctrls.controller('welcome', function($scope) {
    
});

ctrls.controller('browseCtrl', function($scope, $location, $routeParams, ahimsaRestService) {
    ahimsaRestService.getAllBoards().then(function(result) {
        $scope.boards = angular.forEach(result.data, initBoardSum);
        var gex = /^\/board\/(.*)/

        var viewing = null;
        
        var nameL = $location.path().match(gex)
        if (nameL != null  && nameL.length > 0) {
            viewing = nameL[0];
        }
        $scope.openBoard = function(name) {
            viewing = name;
            if (name === "") {
                $location.path("/nilboard");
            } else {
                console.log("This is my path:", $location)
                $location.path("/board/" + name);
            }
        }
        $scope.isOpen = function(name) {
            if (viewing === name) {
                return true;
            } else {
                return false;
            }
        }
    });
})


// Takes in a board and adds processed fields
function initBoardSum(board) {
    board.urlName = encodeURIComponent(board.name);
    return board;
}
