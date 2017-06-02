function createFeatureDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "featureListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/process/{0}/feature/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/process/{0}/feature/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/remove"},
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

    var editControls = isc.ToolStrip.create({
        // width: "100%",
        // height:24,
        members: [
            isc.LayoutSpacer.create({ width:"*" }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/add.png",
                prompt: "Add record",
                click: function() {
                    if (!processList.getSelectedRecord())
                        return false;

                    featureList.startEditingNew({processKey: processList.getSelectedRecord()._key});
                    return false;
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/remove.png",
                prompt: "Remove selected record",
                click: function () {
                    featureList.removeSelectedData(function () {
                        refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
                    })
                }
            })
        ]
    });

    return isc.ListGrid.create({
        ID: "featureList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: featureListDS,
        autoFetchData: false,
        canEdit: true,
        gridComponents:[editControls, "header", "body"],
        selectionUpdated : function (data) {
            console.log (data);
        },
        editComplete: function() {
            console.log ('editComplete: ' + arguments);
            refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
        }
    });
}
