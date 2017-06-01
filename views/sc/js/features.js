isc.DataSource.create({
    dataFormat:"json",
    allowAdvancedCriteria: true,
    ID:"featureList",
    recordXPath: "body",
    operationBindings:[
        {operationType:"fetch", dataProtocol: "", dataURL: "/api/v.0.0.1/process/{id}/feature/list"},
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
                url = operationBinding.dataURL.format (dsRequest.originalData.projectKey);
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