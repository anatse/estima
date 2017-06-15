'use strict';

isc.Log.setPriority("Log", isc.Log.DEBUG);
isc.Log.logDebug ('test messge');

String.prototype.format = String.prototype.f = function(){
    var args = arguments;
    return this.replace(/\{(\d+)\}/g, function(m,n) {
        return args[n] ? args[n] : m;
    });
};

isc.DataSource.create({
    dataURL: "/api/v.0.0.1/user/projects",
    dataFormat: "json",
    allowAdvancedCriteria: true,
    ID:"userProjectListDS",
    recordXPath: "body",
    operationBindings:[
        {operationType:"fetch", dataProtocol: "", dataURL: "/api/v.0.0.1/user/projects"},
        {operationType:"add", dataProtocol:"postMessage", dataURL: "/api/v.0.0.1/project/create"},
        // {operationType:"remove", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL: "/api/v.0.0.1/project/{0}/stage/remove"},
        {operationType:"update", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL:"/api/v.0.0.1/project/create"}
    ],
    transformRequest: function (dsRequest) {
        switch (dsRequest.operationType) {
            case "add":
            case "update":
            case "remove":
                return JSON.stringify(dsRequest.data, function(key, value) {
                    return value;
                });
                break;

            default:
                return dsRequest.data;
        }
    },
    getDataURL: function (dsRequest) {
        var operationBinding = this.getOperationBinding(dsRequest);
        var url = "";
        switch (dsRequest.operationType) {
            case "add":
            case "fetch":
                url = operationBinding.dataURL; //.format (dsRequest.originalData.projectKey);
                break;

            default:
                url = operationBinding.dataURL; //.format (dsRequest.oldValues.projectKey);
        }

        return url;
    },
    fields:[{
        name:"number",
        title: "Номер проекта"
    }, {
        name:"description",
        title:"Описание",
        validators:[
        ]
    }, {
        name:"name",
        title:"Название",
        validators:[
        ]
    }, {
        name:"status",
        type: "integer",
        title: "Статус"
    }, {
        name:"startDate",
        type: "datetime",
        format: "",
        title: "Дата старта"
    }, {
        name:"endDate",
        type: "datetime",
        title: "Дата окончания"
    }, {
        name: "flag",
        title: "флаг"
    }, {
        name: "_key",
        primaryKey: true,
        hidden: true
    }]
});

isc.DataSource.create({
    dataFormat: "json",
    allowAdvancedCriteria: true,
    ID:"projectUserDS",
    recordXPath: "body",
    operationBindings:[
        {operationType:"fetch", dataProtocol: "", dataURL:  "/api/v.0.0.1/project/{0}/user/list"},
        {operationType:"add", dataProtocol:"postMessage", dataURL: "/api/v.0.0.1/project/{0}/user/add"},
        {operationType:"remove", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL: "/api/v.0.0.1/project/{0}/user/remove"},
        {operationType:"update", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL:"/api/v.0.0.1/project/{0}/user/add"}
    ],
    transformRequest: function (dsRequest) {
        switch (dsRequest.operationType) {
            case "add":
            case "update":
            case "remove":
                return JSON.stringify(dsRequest.data, function(key, value) {
                    return value;
                });
                break;

            default:
                return dsRequest.data;
        }
    },
    getDataURL: function (dsRequest) {
        var operationBinding = this.getOperationBinding(dsRequest);
        var url = "";
        switch (dsRequest.operationType) {
            case "add":
            case "fetch":
                url = operationBinding.dataURL.format (dsRequest.originalData.projectKey);
                break;

            default:
                url = operationBinding.dataURL.format ((dsRequest.oldValues || dsRequest.originalData).projectKey);
        }

        return url;
    },
    fields:[{
        name:"name",
        title:"Логин",
        validators:[
        ]
    }, {
        name: "role",
        title: "Роль"
    }, {
        name: "_key",
        primaryKey: true,
        hidden: true
    }, {
        name: "displayName",
        title: "Имя пользователя"
    }, {
        name: "projectKey",
        foreignKey: "userProjectListDS._key",
        hidden: true
    }]
});

isc.DataSource.create({
    cacheAllData: false,
    autoFetchData: false,
    dataFormat: "json",
    allowAdvancedCriteria: false,
    ID: "projectStageListDS",
    recordXPath: "body",
    operationBindings:[
        {operationType:"fetch", dataProtocol: "", dataURL: "/api/v.0.0.1/project/{0}/stage/list"},
        {operationType:"add", dataProtocol:"postMessage", dataURL: "/api/v.0.0.1/project/{0}/stage/add"},
        {operationType:"remove", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL: "/api/v.0.0.1/project/{0}/stage/remove"},
        {operationType:"update", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL:"/api/v.0.0.1/project/{0}/stage/add"}
    ],
    transformRequest: function (dsRequest) {
        switch (dsRequest.operationType) {
            case "add":
            case "update":
            case "remove":
                return JSON.stringify(dsRequest.data, function(key, value) {
                    return value;
                });
                break;

            default:
                return dsRequest.data;
        }
    },
    getDataURL: function (dsRequest) {
        var operationBinding = this.getOperationBinding(dsRequest);
        var url = "";
        switch (dsRequest.operationType) {
            case "add":
            case "fetch":
                url = operationBinding.dataURL.format (dsRequest.originalData.projectKey);
                break;

            default:
                url = operationBinding.dataURL.format (dsRequest.oldValues.projectKey);
        }

        return url;
    },
    fields:[{
        name:"description",
        title:"Описание",
        validators:[
        ]
    }, {
        name:"name",
        title:"Название",
        validators:[
        ]
    }, {
        name:"status",
        type: "integer",
        title: "Статус"
    }, {
        name:"startDate",
        type: "datetime",
        format: "",
        title: "Дата старта"
    }, {
        name:"endDate",
        type: "datetime",
        title: "Дата окончания"
    }, {
        name: "_key",
        primaryKey: true,
        hidden: true
    }, {
        name: "projectKey",
        foreignKey: "userProjectListDS._key",
        hidden: true
    }]
});

isc.DataSource.create({
    dataFormat:"json",
    allowAdvancedCriteria: true,
    ID:"stageProcessListDS",
    recordXPath: "body",
    operationBindings:[
        {operationType:"fetch", dataProtocol: "", dataURL: "/api/v.0.0.1/stage/{0}/process/list"},
        {operationType:"add", dataProtocol:"postMessage", dataURL: "/api/v.0.0.1/stage/{0}/process/create"},
        {operationType:"remove", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL: "/api/v.0.0.1/process/{0}/remove"},
        {operationType:"update", dataProtocol:"postMessage", requestProperties:{httpMethod:"POST"}, dataURL:"/api/v.0.0.1/process/{0}"}
    ],
    transformRequest: function (dsRequest) {
        switch (dsRequest.operationType) {
            case "add":
            case "update":
            case "remove":
                return JSON.stringify(dsRequest.data, function(key, value) {
                    return value;
                });
                break;

            default:
                return dsRequest.data;
        }
    },
    getDataURL: function (dsRequest) {
        var operationBinding = this.getOperationBinding(dsRequest);
        var url = "";
        switch (dsRequest.operationType) {
            case "add":
            case "fetch":
                url = operationBinding.dataURL.format (dsRequest.originalData.stageKey);
                break;

            default:
                url = operationBinding.dataURL.format (dsRequest.originalData._key);
        }

        console.log (url);
        return url;
    },
    fields:[{
        name:"name",
        title:"Название",
        width: 200,
        validators:[
        ]
    },{
        name:"description",
        width: 300,
        title:"Описание",
        validators:[
        ]
    },  {
        name:"status",
        type: "integer",
        title: "Статус"
    }, {
        name: "_key",
        primaryKey: true,
        hidden: true
    }, {
        name: "stageKey",
        foreignKey: "projectStageListDS._key",
        hidden: true
    }]
});

isc.DataSource.create({
    dataFormat:"json",
    allowAdvancedCriteria: true,
    dataURL: "/api/v.0.0.1/users/search",
    ID:"userSearch",
    recordXPath: "body",
    fields:[{
        name:"name",
        title:"Название",
    },{
        name:"displayName",
        title:"Описание",
    },  {
        name:"roles",
        title: "Роли"
    }, {
        name: "_key",
        primaryKey: true,
        hidden: true
    }]
});