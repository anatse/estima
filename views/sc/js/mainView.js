// Create layout
function createSplitPane () {
    let detailPane = isc.DetailViewer.create({
        height: "25%",
        showHeader: true,
        dataSource: "userProjectListDS",
        autoDraw: false
    });

    var splitPane = isc.SectionStack.create ({
        width: "30%",
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
                title: "Детально",
                items: [
                    detailPane
                ]
            }

        ]
    });

    return splitPane;
}

function createTabs () {
    return isc.TabSet.create({
        ID: "topTabSet",
        autoDraw: false,
        tabBarPosition: "top",
        tabs: [
            {title: "Blue"},
            {title: "Green"}
        ]
    });
}

function refreshRelatedGrid (data, parent, child) {
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
        }
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
            console.log (data);
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

isc.HLayout.create({
    width: "100%",
    height: "100%",
    members: [
        createSplitPane(),
        isc.SectionStack.create ({
            width: "20%",
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
                title: "Процессы",
                items: [
                    createProcessGrid()
                ],
                controls:[
                    addProcessButton,
                    removeProcessButton
                ]
            }]
        }),
        createTabs()
    ]
});
