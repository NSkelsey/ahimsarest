var ctrls = angular.module('ombWebAppControllers', ['btford.markdown']);

ctrls.controller('board', function($scope, $routeParams, ahimsaRestService) {
    ahimsaRestService.getBoard($routeParams.board).then(function(result) {
        var board = result.data;
        $scope.board = board;
    });

    $scope.switcheroo = true;

    $scope.id = function(front, bltn) {
        return front + bltn.txid;
    }

    var base = "/static/images/"
    $scope.depthImg = function(bltn) {
        var curHeight = ahimsaRestService.getBlockCount();
        
        if (!angular.isDefined(bltn.blk)) {
            // The bltn is not mined
            return base + "0conf.png"       
        } else {
            // The bltn is in some block
            var diff = curHeight - bltn.blkHeight;

            if (diff > 3) {
                // The bltn is somewhere in the chain
                return base + "totalconf.png"
            }
            // The bltn is less than 5 blocks deep
            return base + (diff + 1) + "conf.png"
        }
    }
});

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
        if (nameL != null  && nameL.length > 1) {
            viewing = nameL[1];
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
