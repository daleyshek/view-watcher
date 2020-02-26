(function(){
    "use strict"
    var report = {
        startTime:null,
        token:null,
        init:function(){
            this.startTime = (new Date()).valueOf();
            this.report();
        },
        report:function(){
            var self = this;
            let http =  new XMLHttpRequest();
            http.open("GET", "{{.HostURL}}/watcher/report")
            if(this.token){
                http.setRequestHeader('Access-Token',this.token);
            }
            http.onload = function(e) {
                if(self.Token == null){
                    self.token = http.getResponseHeader("Access-Token")
                }
              }
            // http.withCredentials = true;
            http.send();
            setTimeout(function(){
                self.report()
            },1000);
        },
    }
    report.init();
})()