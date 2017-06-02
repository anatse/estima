function createFeatureDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "featureListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/process/{0}/feature/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/process/{0}/feature/add"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/process/{0}/remove"},
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                case "add":
                case "update":
                case "remove":
                    return JSON.stringify(dsRequest.data, function (key, value) {
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
                    url = operationBinding.dataURL.format(dsRequest.originalData.processKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
            }

            return url;
        },
        fields: [{
            name: "name",
            title: "Название",
            width: 200,
            validators: []
        }, {
            name: "description",
            width: 300,
            title: "Описание",
            validators: []
        }, {
            name: "status",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "processKey",
            foreignKey: "stageProcessListDS._key",
            hidden: true
        }]
    });
}

function createFeatureGrid () {
    createFeatureDS ();

    return isc.ListGrid.create({
        ID: "featureList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: featureListDS,
        autoFetchData: false,
        canEdit: true,
        selectionUpdated : function (data) {
            console.log (data);
        }
    });
}
