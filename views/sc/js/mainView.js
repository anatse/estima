// Create layout
function createSplitPane () {
    var detailPane = isc.DetailViewer.create({
        height: "25%",
        showHeader: true,
        dataSource: "userProjectListDS",
        autoDraw: false
    });

    return isc.SectionStack.create ({
        showResizeBar: true,
        visibilityMode: "multiple",
        sections: [{
                expanded: true,
                title: "Проекты",
                items: [
                    createProjectGrid(detailPane)
                ],
                controls:[
                    addProjectButton
                ]
            }, {
                expanded: true,
                title: "Дополнительно",
                items: [
                    isc.TabSet.create({
                        ID: "projectDetailsTab",
                        autoDraw: true,
                        tabBarPosition: "top",
                        tabs: [{
                                title: "Детально",
                                pane: detailPane
                            }, {
                                title: "Пользователи",
                                pane: createUsersGrid()
                            }
                        ]
                    })
                ]
            }

        ]
    });
}

function createTabs () {
    return createFeatureGrid();
}

function refreshRelatedGrid (data, parent, child) {
    child.deselectRecords (child.getSelectedRecords());
    child.invalidateCache();
    if (!data) {
        child.fetchData({"key": "-1"})
    } else {
        child.fetchRelatedData(data, parent.dataSource);
    }
}

function createProjectGrid (detailPane) {
    return isc.ListGrid.create({
        ID: "userProjectList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: userProjectListDS,
        autoFetchData: true,
        dateTimeFormatter: 'MM.YYYY HH:MM',
        canEdit: true,
        fields:[{
                name:"number",
                title:"Номер",
                width: 100
            }, {
                name: "name",
                title: "Название",
                width: 300
            }, {
                name:"startDate",
                title:"Начало",
                type:"datetime",
                width: 100
            }, {
                name: "endDate",
                title: "Окончание",
                type: "datetime",
                width: 100
            }
        ],
        selectionUpdated : function (data) {
            detailPane.setData(data);
            refreshRelatedGrid (data, this, stageList);
            refreshRelatedGrid (data, this, userList);
        },
        editComplete: function() {
            this.fetchData()
        }
    });
}

function createUsersGrid () {
    isc.ToolStrip.create({
        ID: "usersEditControls",
        width: "100%", height:24,
        members: [
            isc.LayoutSpacer.create({ width:"*" }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/add.png",
                prompt: "Add record",
                click: function() {
                    if (!userProjectList.getSelectedRecord())
                        return false;

                    userEditWindow.show();
                    return false;
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/remove.png",
                prompt: "Remove selected record",
                click: function () {
                    userList.removeSelectedData(function () {
                        refreshRelatedGrid(userProjectList.getSelectedRecord(), userProjectList, userList);
                    })
                }
            })
        ]
    });

    createUserWindow ();

    return isc.ListGrid.create({
        ID: "userList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: projectUserDS,
        autoFetchData: false,
        showResizeBar: true,
        canEdit: false,
        gridComponents:[usersEditControls, "header", "body"]
    });
}

function createUserWindow () {
    isc.Window.create({
        ID: "userEditWindow",
        title: "User Edit Window",
        autoSize:true,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        show: function (values) {
            userEditForm.setValues (values);
            this.Super("show", arguments)
        },
        items: [
            isc.DynamicForm.create({
                ID: "userEditForm",
                autoFetchData: false,
                autoDraw: false,
                height: 48,
                padding:4,
                fields: [{
                        type:"header",
                        defaultValue:"User add or edit"
                    }, {
                        name: "_key",
                        title: "Ф.И.О.",
                        type: "comboBox",
                        valueField: "_key",
                        displayField: "displayName",
                        addUnknownValues: false,
                        pickListCellWidth: 350,
                        optionDataSource: "userSearch",
                        filterLocally: false,
                        autoFetchData: true,
                        useClientFiltering: false,
                        length: 255,
                        required: true
                    },{
                        name: "role",
                        title: "Роль",
                        type: "select",
                        valueMap: {
                            "OWNER" : "Владелец всего",
                            "PO" : "Владелец продукта",
                            "RTE" : "Сметчик",
                            "ARCHITECT" : "Архитектор",
                            "BA": "Бизнес аналитик",
                            "SA": "Системный аналитик",
                            "PM" : "ТехРук",
                            "TM" : "Разработчик"
                        },
                        required: true
                    }
                ]
            }), isc.HLayout.create ({
                width: "100%",
                height: 10,
                members:[
                    isc.Button.create({
                        name: "validateBtn",
                        title: "Сохранить",
                        click: function () {
                            if (userEditForm.validate()) {
                                var values = userEditForm.getValues();
                                values.projectKey = userProjectList.getSelectedRecord()._key;
                                userList.dataSource.addData(values, function(dsResponse, data, dsRequest) {
                                    var resp = JSON.parse(dsResponse.httpResponseText);
                                    if (!resp.success) {
                                        alert (dsResponse.error);
                                    } else {
                                        console.log(dsResponse);
                                        refreshRelatedGrid(userProjectList.getSelectedRecord(), userProjectList, userList);
                                        userEditWindow.close();
                                    }
                                });


                            }
                        }
                    }), isc.Button.create({
                        name: "hideBtn",
                        title: "Отмена",
                        click: "userEditWindow.hide()"
                    })
                ]
            })
        ]
    });
}

function createStagesGrid () {
    return isc.ListGrid.create({
        ID: "stageList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: projectStageListDS,
        autoFetchData: false,
        showResizeBar: true,
        canEdit: true,
        fields:[{
            name: "name",
            title: "Название",
            width: 300
        }, {
            name:"startDate",
            title:"Начало",
            type:"datetime",
            width: 100
        }, {
            name: "endDate",
            title: "Окончание",
            type: "datetime",
            width: 100
        }],
        selectionUpdated : function (data) {
            refreshRelatedGrid (data, this, processList);
        },
        editComplete: function() {
            refreshRelatedGrid(userProjectList.getSelectedRecord(), userProjectList, stageList);
        }
    });
}

function createProcessGrid () {
    return isc.ListGrid.create({
        ID: "processList",
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: stageProcessListDS,
        autoFetchData: false,
        canEdit: true,
        selectionUpdated : function (data) {
            refreshRelatedGrid (data, this, featureList);
        },
        editComplete: function() {
            refreshRelatedGrid(stageList.getSelectedRecord(), stageList, processList);
        }
    });
}

function createButtons () {
    // Project buttons
    isc.ImgButton.create({
        ID: "addProjectButton",
        autoDraw: false,
        src: "[SKIN]actions/add.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "userProjectList.startEditingNew();return false;"
    });

    isc.ImgButton.create({
        ID: "removeProjectButton",
        autoDraw: false,
        src: "[SKIN]actions/remove.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "userProjectList.removeSelectedData();return false;"
    });

    // Stage buttons
    isc.ImgButton.create({
        ID: "addStageButton",
        autoDraw: false,
        src: "[SKIN]actions/add.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "stageList.startEditingNew({projectKey: userProjectList.getSelectedRecord()._key});return false;"
    });

    isc.ImgButton.create({
        ID: "removeStageButton",
        autoDraw: false,
        src: "[SKIN]actions/remove.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "stageList.removeSelectedData();return false;"
    });

    // Process buttons
    isc.ImgButton.create({
        ID: "addProcessButton",
        autoDraw: false,
        src: "[SKIN]actions/add.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "processList.startEditingNew({stageKey: stageList.getSelectedRecord()._key});return false;"
    });

    isc.ImgButton.create({
        ID: "removeProcessButton",
        autoDraw: false,
        src: "[SKIN]actions/remove.png", size: 16,
        showFocused: false, showRollOver: false, showDown: false,
        click: "processList.removeSelectedData();return false;"
    });
}

createButtons ();

isc.VLayout.create ({
    width: "100%",
    height: "100%",
    members: [
        isc.ToolStrip.create({
            width: "100%",
            height: 24,
            members: [
                isc.ToolStripButton.create({
                    icon: "[SKIN]/actions/configure.png",
                    title: "Рассчитать стоимость проекта"
                }),
                isc.ToolStripButton.create({
                    icon: "[SKIN]/actions/download.png",
                    title: "Експорт в CSV"
                })
            ]
        }),
        isc.HLayout.create({
            members: [
                createSplitPane(),
                isc.SectionStack.create ({
                    showResizeBar: true,
                    width: "15%",
                    visibilityMode: "multiple",
                    sections: [{
                        expanded: true,
                        title: "Этапы",
                        items: [
                            createStagesGrid()
                        ],
                        controls:[
                            addStageButton,
                            removeStageButton
                        ]
                    }, {
                        expanded: true,
                        title: "E2E Процессы",
                        items: [
                            createProcessGrid()
                        ],
                        controls:[
                            addProcessButton,
                            removeProcessButton
                        ]
                    }]
                }),
                createTabs(),
                createTslayput()
            ]
        })
    ]
});

