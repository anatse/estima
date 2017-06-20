String.prototype.format = String.prototype.f = function(){
    var args = arguments;
    return this.replace(/\{(\d+)\}/g, function(m,n) {
        return args[n] ? args[n] : m;
    });
};

function refreshRelatedGrid (data, parent, child) {
    child.deselectRecords (child.getSelectedRecords());
    child.invalidateCache();
    if (!data) {
        child.fetchData({"key": "-1"})
    } else {
        child.fetchRelatedData(data, parent.dataSource);
    }
}

function createCalcUnitDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "cuListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/cu/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cu/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cu/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cu/{0}/remove"},
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
                case "update":
                case "remove":
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
                    break;

                default:
                    url = operationBinding.dataURL;
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
            name: "version",
            width: 60,
            title: "Версия"
        }, {
            name: "changed",
            type: "datetime",
            title: "Дата изменения"
        }, {
            name: "status",
            type: "integer",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }]
    });
}

function createCalcUnitPriceDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "cuPriceListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/cu/{0}/listPrice"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cu/{0}/addPrice"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cuprice/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/cuprice/{0}/remove"},
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
                case "update":
                case "remove":
                    url = operationBinding.dataURL.format(dsRequest.oldValues._key);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData.calcUnitKey);
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
            name: "storyPoints",
            type: "double",
            width: 60,
            title: "Стоимость"
        }, {
            name: "group",
            width: 160,
            valueMap: {
                "Analyze": "Анализ",
                "Integration": "Интеграция",
                "Development": "Разработка",
                "DevOps": "Поддержка",
                "Architect": "Архитектура",
                "SysArch": "Прикладная архитектура",
                "IntArch": "Интеграционная архитектура"
            },
            title: "Группа"
        }, {
            name: "changed",
            type: "datetime",
            title: "Дата изменения"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "calcUnitKey",
            foreignKey: "cuListDS._key",
            hidden: true
        }]
    });
}

function createCalcUnitGrid () {
    createCalcUnitDS();

    return isc.ListGrid.create({
        ID: "cuList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: cuListDS,
        autoFetchData: true,
        canEdit: true,
        selectionUpdated : function (data) {
            refreshRelatedGrid(cuList.getSelectedRecord(), cuList, cuPriceList);
        },
        editComplete: function() {
            cuList.refreshData();
        }
    });
}

function createCalcUnitPriceGrid () {
    createCalcUnitPriceDS();

    return isc.ListGrid.create({
        ID: "cuPriceList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: cuPriceListDS,
        autoFetchData: false,
        canEdit: true,
        editComplete: function() {
            cuList.refreshData();
        }
    });
}

function createCalcUnitlayput () {
    return isc.SectionStack.create ({
        width: "100%",
        height: "100%",
        showResizeBar: true,
        visibilityMode: "multiple",
        sections: [{
            expanded: true,
            title: "Calc Unit",
            items: [
                createCalcUnitGrid()
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить",
                    click: function() {
                        cuList.startEditingNew();
                        return false;
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/remove.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Удалить",
                    click: function () {
                        tsList.removeSelectedData(function () {
                            cuList.refreshData();
                        })
                    }
                })
            ]
        }, {
            expanded: true,
            title: "Calc Unit Price",
            items: [
                createCalcUnitPriceGrid()
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить комментарий",
                    click: function() {
                        if (!cuList.getSelectedRecord())
                            return false;

                        cuPriceList.startEditingNew({calcUnitKey: cuList.getSelectedRecord()._key});
                        return false;
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/remove.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Удалить",
                    click: function () {
                        cuPriceList.removeSelectedData(function () {
                            refreshRelatedGrid(cuList.getSelectedRecord(), cuList, cuPriceList);
                        })
                    }
                })
            ]
        }]
    });
}

createCalcUnitlayput ();