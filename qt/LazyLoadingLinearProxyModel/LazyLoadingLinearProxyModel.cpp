#include "LazyLoadingLinearProxyModel.h"

namespace Status::Helpers {

LazyLoadingLinearProxyModel::LazyLoadingLinearProxyModel(QAbstractItemModel *sourceModel,
                                                         size_t itemsCount, size_t batchSize, QObject *parent)
    //: QAbstractProxyModel(parent)
    : QIdentityProxyModel(parent)
    , m_itemsCount(itemsCount)
    , m_batchSize{batchSize}
    , m_batchSizeRequested{static_cast<int>(sourceModel->rowCount() == 0 ? -1 : sourceModel->rowCount()/batchSize)}
    , m_waitingForValueRole(std::max_element(sourceModel->roleNames().constKeyValueBegin(),
                                             sourceModel->roleNames().constKeyValueEnd(),
                                             [](const auto &a, const auto &b) {
                                                 return a.first < b.first;
                                             })->first + 1)
{
    setSourceModel(sourceModel);
    assert(sourceModel->roleNames().size() > 0);

//    connect(sourceModel, &QAbstractItemModel::rowsAboutToBeInserted, this, [this](const QModelIndex &parent, int start, int end){
//        int startDataChanged = start < m_itemsCount ? start : -1;
//        int endDataChanged = end < m_itemsCount ? end : m_itemsCount - 1;
//        if(start >= 0)
//            emit dataChanged(createIndex(startDataChanged, 0), createIndex(endDataChanged, 0), roleNames().keys());
//        int startInsertChanged = start >= m_itemsCount ? start : m_itemsCount;
//        int endInsertChanged = end >= m_itemsCount ? end : -1;
//        if(endInsertChanged >= 0) {
//            beginInsertRows(mapFromSource(parent), startInsertChanged, endInsertChanged);
//            endInsertRows();
//        }
//    });
}

QVariant LazyLoadingLinearProxyModel::data(const QModelIndex &index, int role) const
{
    const auto &model = *sourceModel();
    const auto needed = index.row();
    if(needed >= model.rowCount()) {
        const auto neededBatch = needed/m_batchSize;
        for(auto batch = m_batchSizeRequested + 1; batch <= neededBatch; ++batch) {
            m_batchSizeRequested = batch;
            emit doLoadBatch(batch);
        }
        if(role == m_waitingForValueRole)
            return true;
        else
            return QVariant();
    } else {
        if(role == m_waitingForValueRole)
            return false;
        else
            return QAbstractProxyModel::data(index, role);
    }
}

bool LazyLoadingLinearProxyModel::canFetchMore(const QModelIndex &parent) const {
    return sourceModel()->rowCount() < m_itemsCount;
}

void LazyLoadingLinearProxyModel::fetchMore(const QModelIndex &parent)
{
    m_batchSizeRequested = m_batchSizeRequested + 1;
    emit doLoadBatch(m_batchSizeRequested);
}

//int LazyLoadingLinearProxyModel::columnCount(const QModelIndex &parent) const
//{
//    return 1;
//}

//int LazyLoadingLinearProxyModel::rowCount(const QModelIndex &parent) const
//{
//    const auto &model = *sourceModel();
//    return (model.rowCount() >= m_itemsCount) ? model.rowCount() : m_itemsCount;
//}

//QModelIndex LazyLoadingLinearProxyModel::index(int row, int column, const QModelIndex &parent) const
//{
//    return createIndex(row, column);
//}

//QModelIndex LazyLoadingLinearProxyModel::mapFromSource(const QModelIndex &sourceIndex) const
//{
//    return index(sourceIndex.row(), sourceIndex.column());
//}

//QModelIndex LazyLoadingLinearProxyModel::mapToSource(const QModelIndex &proxyIndex) const
//{
//    return sourceModel()->index(proxyIndex.row(), proxyIndex.column());
//}

//QModelIndex LazyLoadingLinearProxyModel::parent(const QModelIndex &index) const { return QModelIndex(); }

QHash<int, QByteArray> LazyLoadingLinearProxyModel::roleNames() const {
    // TODO cache it, it is used in other palces also
    assert(sourceModel() != nullptr);
    auto newRoles{sourceModel()->roleNames()};
    newRoles.insert(m_waitingForValueRole, m_waitingForValueRoleName);
    return newRoles;
}

void LazyLoadingLinearProxyModel::resetRequestedBatch(int batch)
{
    // This works as expected because the increment part can only go up and we failed lower
    if(batch <= m_batchSizeRequested) {
        m_batchSizeRequested = batch - 1;
        int start = batch * m_batchSize;
        int end = start + m_batchSize - 1;
        beginRemoveRows(QModelIndex(), start, end);
        endRemoveRows();
    }
};

}
