package repository

const (
	insertIntoOrder = `
		INSERT INTO "order" (
			order_uid,
			track_number,
			entry,
			locale,
			internal_signature,
			customer_id,
			delivery_service,
			shardkey,
			sm_id,
			date_created,
			oof_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		);`
	insertIntoDelivery = `
		INSERT INTO delivery (
			order_uid,
			name,
			phone,
			zip,
			city,
			address,
			region,
			email
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8 
		);`

	insertIntoPayment = `
		INSERT INTO payment (
			order_uid,
			transaction,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		);`

	insertIntoItem = `
		INSERT INTO item (
			order_uid,
			chrt_id,
			track_number,
			price,
			rid,
			name,
			sale,
			size,
			total_price,
			nm_id,
			brand,
			status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 
		)`
)
