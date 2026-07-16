# Service documentation

## Structure

The service struct consists of references to
repos and database. It should be created by
function NewBonusService from the database reference
to create all repos for this database.

## Function and methods

### NewBonusService

#### Arguments

Database reference

#### Returns

Reference to the new instance of service

#### Dependencies

NewBalanceRepo - storage layer

NewBatchRepo - storage layer

NewHoldRepo - storage layer

NewLedgerRepo - storage layer

NewHoldBatchRepo - storage layer

### Accrue

#### Arguments

Context, reference to ledger entry, number of days

#### Returns

ID of new batch, error

#### Dependencies

validate - service layer

BatchRepo.CreateBatch - storage layer

BalanceRepo.UpdateBalance - storage layer

LedgerRepo.Insert - storage layer

#### Important notes

Given entry have to contain user id, amount of points,
and idempotency key. The method completes
entry with the rest of the data and validates request.

### Hold

#### Arguments

Context, reference to ledger entry,
order id, number of hours

#### Returns

ID of new hold, error

#### Dependencies

validate - service layer

BatchRepo.GetExpiringBatchesForUpdate - storage layer

BatchRepo.DecreaseBatchRemaining - storage layer

HoldBatchRepo.CreateHoldBatch - storage layer

HoldRepo.CreateHold - storage layer

BalanceRepo.UpdateBalance - storage layer

LedgerRepo.Insert - storage layer

#### Important notes

Given entry have to contain user id, amount of points,
and idempotency key. Order id have to be unique
for user. The method completes
entry with the rest of the data and validates request.

### Cancel hold

#### Arguments

Context, reference to ledger entry,
order id

#### Returns

ID of cancelled hold, error

#### Dependencies

validate - service layer

BatchRepo.IncreaseBatchRemaining - storage layer

HoldBatchRepo.GetHoldBatchesByHoldID - storage layer

HoldRepo.GetHoldByOrderID - storage layer

HoldRepo.UpdateHoldStatus - storage layer

BalanceRepo.UpdateBalance - storage layer

LedgerRepo.Insert - storage layer

#### Important notes

Given entry have to contain user id
and idempotency key. Order id have to be unique
for user. The method completes
entry with the rest of the data and validates request.

### ConfirmWithdraw

#### Arguments

Context, reference to ledger entry,
order id

#### Returns

ID of confirmed hold, error

#### Dependencies

validate - service layer

HoldRepo.GetHoldByOrderID - storage layer

HoldRepo.UpdateHoldStatus - storage layer

BalanceRepo.UpdateBalance - storage layer

LedgerRepo.Insert - storage layer

#### Important notes

Given entry have to contain user id
and idempotency key. Order id have to be unique
for user. The method completes
entry with the rest of the data and validates request.

### GetAvailablePoints

#### Arguments

Context, user id

#### Returns

Number of all available points for the user, error

#### Dependencies

BalanceRepo.GetBalance - storage layer

### GetExpiringAvailablePoints

#### Arguments

Context, user id, number of days

#### Returns

Number of all currently available points for the user
which will be expired in the given number
of days, error

#### Dependencies

BatchRepo.GetExpiringSum - storage layer

### GetHeld

#### Arguments

Context, user id

#### Returns

Number of all currently held points for the user,
error

#### Dependencies

BalanceRepo.GetBalance - storage layer

### GetHeldByOrderId

#### Arguments

Context, user id, order id

#### Returns

Number of held for the order points for the user,
error

#### Dependencies

HoldRepo.GetHoldByOrderID - storage layer

### GetHistory

#### Arguments

Context, user id, limit, offset

#### Returns

Slice of ledger entries for user of limited size
ignoring offset of first entries, error

#### Dependencies

LedgerRepo.GetHistory - storage layer

### validate

#### Arguments

Context, user id, number of required points,
idempotency key, arbitrary number of variables which
have to be positive

#### Returns

transaction reference, reference to user's balance,
error

#### Dependencies

BalanceRepo.GetBalanceForUpdate - storage layer

LedgerRepo.GetByExternalKey - storage layer

#### Important notes

May be called only within service methods (it is a
private method). Locks user's balance within the
transaction so it is impossible to run in parallel
two functions using validate in the same time for
the same user within different transaction.