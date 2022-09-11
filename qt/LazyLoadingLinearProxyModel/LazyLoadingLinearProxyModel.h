#pragma once

#include <QIdentityProxyModel>
#include <QQmlEngine>

#include <functional>

namespace Status::Helpers {

/// Helper to have lazy loading of models
/// \todo test loading request math
/// \todo test lazy data loading
class LazyLoadingLinearProxyModel final : public QIdentityProxyModel//QAbstractProxyModel
{
    Q_OBJECT

public:
    using RequestToLoadBatchFn = std::function<void(int)>;

    /// \note The sourceModel owner is responsible keeping a continous structure for the source model when adding batches
    explicit LazyLoadingLinearProxyModel(QAbstractItemModel *sourceModel,
                                         size_t itemsCount,
                                         size_t batchSize,
                                         QObject *parent = nullptr);

    QVariant data(const QModelIndex &index, int role) const override;

    bool canFetchMore(const QModelIndex &parent) const override;
    void fetchMore(const QModelIndex &parent) override;

//    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
//    int rowCount(const QModelIndex &parent = QModelIndex()) const override;

//    QModelIndex index(int row, int column, const QModelIndex &parent = QModelIndex()) const override;

//    QModelIndex mapFromSource(const QModelIndex & sourceIndex = QModelIndex()) const override;

//    QModelIndex mapToSource(const QModelIndex & proxy) const override;

//    QModelIndex parent(const QModelIndex &index) const override;;

    QHash<int, QByteArray> roleNames() const override;

    void resetRequestedBatch(int batch);

signals:

    void doLoadBatch(size_t batch) const;

private:
    size_t m_itemsCount;
    size_t m_batchSize;
    mutable std::atomic<int> m_batchSizeRequested;
    RequestToLoadBatchFn m_loadCallback;

    const QByteArray m_waitingForValueRoleName{"waitingForValue"};
    const int m_waitingForValueRole{};
};

}
