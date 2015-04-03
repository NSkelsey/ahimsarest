'use strict';

var ombFilters = angular.module('ombWebAppFilters', []);

ombFilters.filter('nicedate', function() {
    return function(date) {
        var options = { year: '2-digit', month: '2-digit', day: '2-digit' }
        return date.toLocaleDateString('en-US', options)
    }
});

ombFilters.filter('nicedatetime', function() {
    var options = { 
        year: '2-digit',
        month: '2-digit',
        day: '2-digit',
        hour: 'numeric',  
        minute: 'numeric',
        timeZone: 'UTC',
    };
    var formater = new Intl.DateTimeFormat('en-US', options);
    return function(utcsecs) {
        var d = new Date(utcsecs*1000); 
        console.log(d);
        return formater.format(d);
    }
});

ombFilters.filter('epochdate', function() {
    return function(utcsecs) {
        return new Date(utcsecs*1000)
    }
});

ombFilters.filter('plural', function() {
    return function(num, word) {
        if (num == 1) {
            return num + " " + word.slice(0, word.length-1)
        } else {
            return num + " " + word
        }
    }
});
