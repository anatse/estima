function foo(params) {
    let db = require('internal').db,
        toCol = db._collection (params.toColName),
        edgeCol = db._collection(params.edgeColName),
        doc = db._collection(params.fromColName).document(params.fKey),
        activeText = null,
        numActives = 0;

    db._createStatement({
        query: "FOR v, e IN OUTBOUND '" + doc._id + "' " + params.edgeColName + " FILTER e.label == 'text' && v.active RETURN v"
    }).execute().toArray().forEach((d) => {
        activeText = d;
        numActives++;
    });

    if (numActives > 1)
        throw ("Found more than one active texts for given feature. Please, contact with administrator to fix problem");

    var version = 0
    if (activeText) {
        // Set inactive for previous text
        activeText.active = false;
        toCol.replace(activeText._id, activeText);
        version = activeText.version;
    }

    // Save new text with ++version
    var newText = {
        text: params.text.text,
        active: true,
        version: version + 1
    };
    let toDoc = toCol.save (newText);

    // Add edge between text and object
    edgeCol.save ({_from: doc._id, _to: toDoc._id, label: 'text'});
    return {success: true, entityKey: toDoc._key};
}
