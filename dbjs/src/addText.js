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

    // Set inactive for previous text
    activeText.active = false;
    toCol.save (activeText);

    // Save new text
    params.text.active = true;
    let toDoc = toCol.save (params.text);

    // Add edge between text and object
    let edge = {_from: doc._id, _to: toDoc._id};
    for (let v in p.props) {
        edge[v] = p.props[v];
    }
    edge = edgesCol.save (edge);
    return {success: true, entityKey: toDoc._key};
}
