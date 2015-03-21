angular.module('ahimsaApp.blocksView', ['ngRoute'])
    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/blocks/:day', {
                  controller: 'blockDayCtrl',
                  templateUrl: 'views/blocks/blocks.html'
        });
    }])
    .controller('blockDayCtrl', function($scope, $routeParams, ahimsaRestService) {

        $scope.day = parseDayStr($routeParams.day)

        $scope.$on('back', function() {
            var prevDay = new Date($scope.day)
            prevDay.setDate(prevDay.getDate()-1);
            $location.path(formatDayStr(prevDay));
        });

        $scope.$on('forward', function() {
            var prevDay = new Date($scope.day)
            prevDay.setDate(prevDay.getDate()+1);   
            $location.path(formatDayStr(nextDay));
        });
        
        ahimsaRestService.getBlocksDay($routeParams.day).then(function(result) {
            $scope.blocks = result.data
        });
    })
