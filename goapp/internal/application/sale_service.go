package application

import (
	"errors"
	"fmt"
	"opticav2/internal/domain"
	"time"

	"gorm.io/gorm" // For transaction management if needed at service level
)

type SaleService struct {
	SaleRepo    domain.SaleRepository
	PaymentRepo domain.PaymentRepository
	ProductRepo domain.ProductRepository
	ClientRepo  domain.ClientRepository
	DB          *gorm.DB // Inject DB for transaction management at service level if needed
}

func NewSaleService(
	saleRepo domain.SaleRepository,
	paymentRepo domain.PaymentRepository,
	productRepo domain.ProductRepository,
	clientRepo domain.ClientRepository,
	db *gorm.DB, // Pass GORM DB instance for transactions
) *SaleService {
	return &SaleService{
		SaleRepo:    saleRepo,
		PaymentRepo: paymentRepo,
		ProductRepo: productRepo,
		ClientRepo:  clientRepo,
		DB:          db,
	}
}

func (s *SaleService) CreateSale(req domain.CreateSaleRequest, userID uint) (*domain.Sale, error) {
	// 1. Validate ClientID
	_, err := s.ClientRepo.GetByID(int(req.ClientID)) // Assuming Client ID is int
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			return nil, fmt.Errorf("client with ID %d not found: %w", req.ClientID, domain.ErrInvalidSaleData)
		}
		return nil, fmt.Errorf("error validating client: %w", err)
	}

	// 2. Prepare SaleItems and calculate TotalAmount
	var saleItems []domain.SaleItem
	var totalAmount float64 = 0

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("sale must have at least one item: %w", domain.ErrInvalidSaleData)
	}

	for _, itemReq := range req.Items {
		product, err := s.ProductRepo.GetByID(int(itemReq.ProductID)) // Assuming Product ID is int
		if err != nil {
			if errors.Is(err, domain.ErrProductNotFound) {
				return nil, fmt.Errorf("product with ID %d not found for sale item: %w", itemReq.ProductID, domain.ErrInvalidSaleData)
			}
			return nil, fmt.Errorf("error fetching product %d: %w", itemReq.ProductID, err)
		}

		unitPrice := product.Price // Default to product's current price
		if itemReq.UnitPrice > 0 { // Allow overriding price at time of sale
			unitPrice = itemReq.UnitPrice
		}

		itemTotalPrice := unitPrice * float64(itemReq.Quantity)
		totalAmount += itemTotalPrice

		saleItems = append(saleItems, domain.SaleItem{
			ProductID:          itemReq.ProductID,
			Quantity:           itemReq.Quantity,
			UnitPrice:          unitPrice,
			TotalPrice:         itemTotalPrice,
			ProductDescription: product.Description, // Denormalization
			ProductCode:        product.Code,        // Denormalization
		})
	}

	// 3. Calculate final amounts and status
	finalAmount := totalAmount - req.DiscountAmount
	if finalAmount < 0 {
		return nil, fmt.Errorf("final amount cannot be negative after discount: %w", domain.ErrInvalidSaleData)
	}

	amountPaid := req.InitialPayment
	balanceDue := finalAmount - amountPaid

	if amountPaid > finalAmount {
		// This could be an overpayment scenario, or an error depending on business rules
		// For now, let's assume initial payment cannot exceed final amount.
		// Or, it could mean change is due if payment method allows.
		// Simplification: initial payment is capped at final_amount or is exact.
		// If initial_payment > final_amount, it could be an error or change is due.
		// Let's assume it's an error for now if it's significantly larger.
		// For simplicity, we'll allow it, but balanceDue will be negative (meaning change or credit)
		// This needs clarification based on business rules.
		// Alternative: cap amountPaid at finalAmount for initial payment if no change is handled.
		// amountPaid = math.Min(req.InitialPayment, finalAmount)
		// balanceDue = finalAmount - amountPaid
	}


	status := "Pending" // Default status
	if balanceDue <= 0 {
		status = "Completed"
	} else if amountPaid > 0 && balanceDue > 0 {
		status = "Partial"
	}


	// 4. Prepare Sale object
	saleDate := time.Now()
	if req.SaleDate != "" {
		parsedDate, err := time.Parse("2006-01-02 15:04:05", req.SaleDate) // Or just YYYY-MM-DD
		if err == nil {
			saleDate = parsedDate
		} else {
            parsedDate, err = time.Parse("2006-01-02", req.SaleDate)
            if err == nil {
                saleDate = parsedDate
            }
        }
	}

	sale := &domain.Sale{
		ClientID:         req.ClientID,
		UserID:           userID,
		SaleDate:         saleDate,
		TotalAmount:      totalAmount,
		DiscountAmount:   req.DiscountAmount,
		FinalAmount:      finalAmount,
		AmountPaid:       amountPaid, // Initially, this is the initial payment
		BalanceDue:       balanceDue,
		PaymentMethodID:  req.PaymentMethodID, // Initial payment method
		Status:           status,
		PrescriptionData: req.PrescriptionData,
		Observations:     req.Observations,
	}

	// 5. Prepare initial Payment object (if any)
	var initialPaymentRecord *domain.Payment
	if req.InitialPayment > 0 {
		initialPaymentRecord = &domain.Payment{
			// SaleID will be set by the repository after Sale is created
			PaymentDate:       sale.SaleDate, // Typically same as sale date for initial payment
			Amount:            req.InitialPayment,
			PaymentMethodID:   req.PaymentMethodID,
			ProcessedByUserID: userID,
		}
	}

	// 6. Call Repository Create (handles transaction for sale, items, payment, stock)
	err = s.SaleRepo.Create(sale, saleItems, initialPaymentRecord)
	if err != nil {
		return nil, err // Error from repository (includes stock errors)
	}

	// The 'sale' object should now have its ID and potentially nested items/payments if GORM handles it.
	// The repo's Create method should ensure 'sale.ID' is populated.
	// For returning the full sale details, a GetByID might be needed if repo doesn't backfill associations.
	// For now, assume 'sale' object from repo.Create is sufficiently populated (at least with ID).
	return s.SaleRepo.GetByID(sale.ID) // Fetch fresh record with all associations
}

func (s *SaleService) GetSale(saleID uint, _ uint) (*domain.Sale, error) {
	// TODO: Add authorization check: userID (second param) against sale.UserID or user role
	sale, err := s.SaleRepo.GetByID(saleID)
	if err != nil {
		return nil, err // Handles domain.ErrSaleNotFound from repo
	}
	return sale, nil
}

func (s *SaleService) ListSales(userID uint, filters map[string]interface{}) ([]domain.Sale, error) {
	// TODO: Implement proper filtering logic, possibly based on user role.
	// If user is not admin, restrict filters to their own sales:
	// if !isUserAdmin(userID) { // isUserAdmin would check user role
	//    filters["user_id"] = userID
	// }
	return s.SaleRepo.GetAll(filters)
}

func (s *SaleService) AddPaymentToSale(saleID uint, req domain.AddPaymentRequest, processedByUserID uint) (*domain.Sale, error) {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Use a transactional version of repositories if they are not already transaction-aware
        // For simplicity, assume repositories here can work with `tx` if needed,
        // or methods are called that don't need explicit tx passing because they are simple ops.
        // Better: saleRepoTx := s.SaleRepo.WithTx(tx) etc.

        sale, err := s.SaleRepo.GetByID(saleID) // Fetch with original DB context or tx context
        if err != nil {
            return err
        }

        if sale.Status == "Completed" || sale.Status == "Cancelled" {
            return fmt.Errorf("cannot add payment to sale with status: %s", sale.Status)
        }

        if req.Amount <= 0 {
            return errors.New("payment amount must be positive")
        }
        // Potential check: if req.Amount > sale.BalanceDue (handle overpayment?)

        paymentDate := time.Now()
        if req.PaymentDate != "" {
            parsedDate, err := time.Parse("2006-01-02 15:04:05", req.PaymentDate)
            if err == nil {
                paymentDate = parsedDate
            } else {
                 parsedDate, err = time.Parse("2006-01-02", req.PaymentDate)
                 if err == nil {
                    paymentDate = parsedDate
                }
            }
        }

        payment := &domain.Payment{
            SaleID:            saleID,
            PaymentDate:       paymentDate,
            Amount:            req.Amount,
            PaymentMethodID:   req.PaymentMethodID,
            ProcessedByUserID: processedByUserID,
        }

        if err := s.PaymentRepo.Create(payment); err != nil { // Assume PaymentRepo.Create uses its own DB instance or is tx-aware
            return fmt.Errorf("failed to create payment record: %w", err)
        }

        // Recalculate sale's payment fields
        sale.AmountPaid += req.Amount
        sale.BalanceDue = sale.FinalAmount - sale.AmountPaid // Ensure FinalAmount is correct

        if sale.BalanceDue <= 0 {
            sale.Status = "Completed"
            // If balanceDue is negative, it implies overpayment/change. Business rule needed.
        } else {
            sale.Status = "Partial" // If previously "Pending" or still "Partial"
        }

        // Update sale within the transaction
        if err := s.SaleRepo.Update(sale); err != nil { // Assume SaleRepo.Update uses its own DB or is tx-aware
            return fmt.Errorf("failed to update sale after payment: %w", err)
        }
        return nil // Commit transaction
    }).Error // Return the error from the transaction block
    // After successful transaction, fetch and return the updated sale
    if err != nil {
        return nil, err
    }
    return s.SaleRepo.GetByID(saleID)
}


func (s *SaleService) CancelSale(saleID uint, _ uint) (*domain.Sale, error) {
    // The `cancelledByUserID` (second param) would be used for logging or authorization checks.

    // For cancelling a sale, especially with stock restoration, it's critical to use a transaction.
    // The SaleRepository's Create method already handles stock *decrement* transactionally with sale items.
    // A symmetric operation is needed for cancellation.
    // Option 1: Service layer manages the transaction.
    // Option 2: Repository has a specific `CancelSale(saleID uint) error` method that handles it all.

    // Let's proceed with Option 1: Service layer transaction.
    err := s.DB.Transaction(func(tx *gorm.DB) error {
        // It's better if repository methods can accept a `*gorm.DB` to participate in caller's transaction.
        // e.g., s.SaleRepo.GetByIDWithTx(tx, saleID), s.ProductRepo.UpdateWithTx(tx, product)
        // Assuming current repo methods use their own DB instance, this transaction won't cover their internal operations.
        // This is a common architectural challenge.
        // For this example, let's assume we modify repo methods or use direct tx operations for critical parts.

        sale, err := s.SaleRepo.GetByID(saleID) // Fetch outside tx or make GetByID tx-aware
        if err != nil {
            return err
        }

        if sale.Status == "Cancelled" {
            return errors.New("sale is already cancelled")
        }
        // Add other checks, e.g., if sale is too old to be cancelled based on policy.

        // Restore stock for each item
        for _, item := range sale.SaleItems {
            // This should ideally use ProductRepo method that accepts `tx`
            // err := s.ProductRepo.IncreaseStock(tx, item.ProductID, item.Quantity)
            // Direct update with `tx` for now:
            if errStock := tx.Model(&domain.Product{}).Where("codproducto = ?", item.ProductID).
                UpdateColumn("existencia", gorm.Expr("existencia + ?", item.Quantity)).Error; errStock != nil {
                return fmt.Errorf("failed to restore stock for product ID %d: %w", item.ProductID, errStock)
            }
        }

        sale.Status = "Cancelled"
        // sale.BalanceDue = sale.FinalAmount - sale.AmountPaid // Recalculate if refunds are part of this, complex.
        // For now, assume cancellation just marks as cancelled and restores stock. Refunds are separate.

        // Update sale status using `tx`
        // err = s.SaleRepo.UpdateWithTx(tx, sale)
        // Direct update with `tx` for now:
        if errUpdate := tx.Save(sale).Error; errUpdate != nil {
            return fmt.Errorf("failed to update sale status to cancelled: %w", errUpdate)
        }

        // TODO: Handle refunds or voiding payments if necessary. This is complex.
        // For now, payments remain as they were, sale is just marked "Cancelled".

        return nil // Commit transaction
    })

    if err != nil {
        return nil, err
    }

    return s.SaleRepo.GetByID(saleID) // Return updated sale
}
