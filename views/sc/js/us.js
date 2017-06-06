
function createUSTextDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "usTextDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/feature/{0}/userstory/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/userstory/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/remove"},
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
                    url = operationBinding.dataURL.format(dsRequest.originalData.featureKey);
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
            name: "who",
            width: 300,
            title: "Кто",
            validators: []
        },{
            name: "what",
            width: 300,
            title: "Что",
            validators: []
        },{
            name: "why",
            width: 300,
            title: "Зачем",
            validators: []
        },{
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
            name: "serial",
            title: "Номер",
            type: "integer",
            primaryKey: true
        }, {
            name: "featureKey",
            foreignKey: "featureListDS._key",
            hidden: true
        }]
    });
}

function createUsGrid () {
    createUSDS();

    var editControls = isc.ToolStrip.create({
        members: [
            isc.LayoutSpacer.create({ width:"*" }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/add.png",
                prompt: "Add record",
                click: function() {
                    if (!featureList.getSelectedRecord())
                        return false;

                    usList.startEditingNew({featureKey: featureList.getSelectedRecord()._key, serial: usList.getData().totalRows + 1});
                    return false;
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/remove.png",
                prompt: "Remove selected record",
                click: function () {
                    usList.removeSelectedData(function () {
                        refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
                    })
                }
            })
        ]
    });

    return isc.ListGrid.create({
        ID: "usList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: usListDS,
        autoFetchData: false,
        canEdit: true,
        gridComponents:[editControls, "header", "body"],
        selectionUpdated : function (data) {
            console.log (data);
        },
        editComplete: function() {
            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
        }
    });
}
