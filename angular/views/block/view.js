
angular.module('ahimsaApp.blockView', ['ngRoute'])
    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/block/:hash', {
            templateUrl: 'views/block/block.html',
            controller: 'blockCtrl',
        });
    }])
    .controller('blockCtrl', function($scope, $routeParams, ahimsaRestService) {
        $scope.$on('Reverse', function () {
            $scope.bltns.reverse()
        });
        
        ahimsaRestService.getBlock($routeParams.hash).then(function(result) {
            $scope.block = result.data
            $scope.bltns = result.data.bltns
        })
    });
