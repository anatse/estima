function foo(p) {
    let db = require('internal').db,
        toCol = db._collection(p.toColName),
        fromCol = db._collection(p.fromColName),
        toDoc, fromDoc;

    if (!p.props['label'])
        throw ('Error creation edge - label is not defined');

    fromDoc = fromCol.document(p.fromId);

    try {
        toDoc = (p.doc._key !== "") ? toDoc = toCol.save(p.doc) : toCol.document(p.doc._key);
    } catch(e) {
        if (e.errorNum === 1202)
            toDoc = toCol.save(p.doc);
        else
            throw (e);
    }

    let edgesCol = db._collection (p.edgeColName);
    let edge = {_from: fromDoc._id, _to: toDoc._id};
    for (let v in p.props) {
        edge[v] = p.props[v];
    }

    edge = edgesCol.save (edge);
    return {success: true, entityKey: toDoc._key};
}