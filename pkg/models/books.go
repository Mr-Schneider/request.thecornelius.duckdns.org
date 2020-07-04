package models

import (
	"database/sql"
	"log"
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/forms"
)

// BooksDB holds the db connection
type BooksDB struct {
	*sql.DB
}

// GetBook retrives a book from the db
func (db *DB) GetBook(id string) (*Book, error) {
	// Query statement
	stmt := `SELECT id, volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created FROM books WHERE volumeid = $1`

	// Execute query
	row := db.QueryRow(stmt, id)
	b := &Book{}

	// Pull data into request
	err := row.Scan(&b.ID, &b.VolumeID, &b.Title, &b.Subtitle, &b.Publisher, &b.PublishedDate, &b.PageCount,
		&b.MaturityRating, &b.Authors, &b.Categories, &b.Description, &b.Uploader, &b.Price, &b.ISBN10, &b.ISBN13,
		&b.ImageLink, &b.Downloads, &b.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

// LatestBooks grabs the latest 10 valid books
func (db *DB) LatestBooks(limit int) (Books, error) {
	// Query statement
	stmt := `SELECT id, volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created FROM books ORDER BY created DESC LIMIT $1`

	// Execute query
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := Books{}

	// Get all the matching requets
	for rows.Next() {
		b := &Book{}

		// Pull data into request
		err := rows.Scan(&b.ID, &b.VolumeID, &b.Title, &b.Subtitle, &b.Publisher, &b.PublishedDate, &b.PageCount,
			&b.MaturityRating, &b.Authors, &b.Categories, &b.Description, &b.Uploader, &b.Price, &b.ISBN10, &b.ISBN13,
			&b.ImageLink, &b.Downloads, &b.Created)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// InsertBooks adds a new book to the library
func (db *DB) InsertBook(new_book *forms.NewBook) (int, error) {
	// Save stored request
	var bookid int

	// Query statement
	stmt := `INSERT INTO books (volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 0, timezone('utc', now())) RETURNING id`

	err := db.QueryRow(stmt, new_book.VolumeID, new_book.Title, new_book.Subtitle, new_book.Publisher, new_book.PublishedDate, new_book.PageCount,
		new_book.MaturityRating, new_book.Authors, new_book.Categories, new_book.Description, new_book.Uploader, new_book.Price, new_book.ISBN10, new_book.ISBN13,
		new_book.ImageLink).Scan(&bookid)
	if err != nil {
		return 0, err
	}

	log.Printf("New book %s uploaded by %s", new_book.Title, new_book.Uploader)

	return bookid, nil
}

// DownloadBook incriments the download counter
func (db *DB) DownloadBook(book_id string, downloads int) {

	// Query statement
	stmt := `UPDATE books SET downloads = $1 WHERE volumeid = $2`

	db.QueryRow(stmt, downloads, book_id)
}

// UpdateBook edits a book attributes
func (db *DB) UpdateBook(book *forms.NewBook) (int, error) {
	// Save stored request
	var bookid int

	// Query statement
	stmt := `UPDATE books SET title = $1, subtitle = $2, publisher = $3, publisheddate = $4, pagecount = $5,
		maturityrating = $6, authors = $7, categories = $8, description = $9, price = $10, isbn10 = $11, isbn13 = $12, imagelink = $13 WHERE volumeid = $14`

	db.QueryRow(stmt, book.Title, book.Subtitle, book.Publisher, book.PublishedDate, book.PageCount,
		book.MaturityRating, book.Authors, book.Categories, book.Description, book.Price, book.ISBN10, book.ISBN13,
		book.ImageLink, book.VolumeID)

	log.Printf("Book %s edited", book.Title)

	return bookid, nil
}