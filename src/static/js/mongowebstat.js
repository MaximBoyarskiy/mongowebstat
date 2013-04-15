mongowebstat = angular.module('mongowebstat', ['statService', 'nodeService']);

function StatListCtrl($scope, $timeout, Stats, Nodes) {
    $scope.stats = [];
    $scope.nodes = Nodes.query();
    $scope.isActive = true;
    $scope.fields = ['node', 'insert', 'query', 'update', 'delete',
            'get more', 'command', 'flushes', 'mapped','virtual', 'resident', 'non mapped',
            'faults' ,'idx miss %', 'qr|qw', 'ar|aw', 'net in|out', 'conn'];

    $scope.filterNet = function(diff) {
        var unit = "b";
        var div = 1000;
        if ( diff >= div ) {
            unit = "k";
            diff /= div;
        };

        if ( diff >= div ) {
            unit = "m";
            diff /= div;
        };

        if ( diff >= div ) {
            unit = "g";
            diff /= div;
        };
        if (diff > 100) {
            diff = diff.toFixed(0)
        } else {
            diff = diff.toFixed(1)
        };

        return diff + unit;
    };

    $scope.filterMem = function(sz) {
        var unit = "m";
        if ( sz > 1024 ) {
            unit = "g";
            sz /= 1024;
        };

        if (sz > 100) {
            sz = sz.toFixed(0)
        } else {
            sz = sz.toFixed(1)
        };

        return sz + unit;
    };

    (function tick() {
        if ($scope.isActive){
            $scope.stats = Stats.query(function(){
                $timeout(tick, 1000);
            });
        }
        else
        {
            $timeout(tick, 1000);
        }

    })();
};




angular.module('statService', ['ngResource']).
    factory('Stats', function($resource){
        return $resource('/stats', {}, {
    query: {method:'GET', params:{}, isArray:false}
  });
});

angular.module('nodeService', ['ngResource']).
    factory('Nodes', function($resource){
        return $resource('/nodes', {}, {
    query: {method:'GET', params:{}, isArray:true}
  });
});

$(document).ready(function () {
    $("td").tooltip({
        'placement': 'right'
    });
});
