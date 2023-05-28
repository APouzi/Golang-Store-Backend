CREATE TABLE IF NOT EXISTS tblProducts (
  ProductID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ProductName VARCHAR(255) NOT NULL,
  ProductDescription TEXT,
  ProductPrice DECIMAL(10,2) NOT NULL,
  SKU VARCHAR(50),
  UPC VARCHAR(50),
  PRIMARY_IMAGE VARCHAR(255) NULL,
  ProductDateAdded TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  ModifiedDate DATETIME NULL
  
);

CREATE TABLE IF NOT EXISTS tblCategoriesPrime (
    CategoryID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesSub (
    CategoryID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesFinal (
    CategoryID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);

CREATE TABLE IF NOT EXISTS tblCatPrimeSub (
  CatPrimeID INT NOT NULL,
  CatSubID INT NOT NULL,
  FOREIGN KEY (CatPrimeID) REFERENCES tblCategoriesPrime (CategoryID) ON DELETE CASCADE,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (CategoryID) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS tblCatSubFinal (
  CatSubID INT NOT NULL,
  CatFinalID INT NOT NULL,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (CategoryID) ON DELETE CASCADE,
  FOREIGN KEY (CatFinalID) REFERENCES tblCategoriesFinal (CategoryID) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS tblCatFinalProd (
  CatFinalID INT NOT NULL,
  ProductID INT NOT NULL,
  FOREIGN KEY (CategoryID) REFERENCES tblCategoriesFinal (CategoryID) ON DELETE CASCADE,
  FOREIGN KEY (CategoryID) REFERENCES tblProducts (CategoryID) ON DELETE CASCADE
  
)

-- Lines 35 to 58 are all Denormalization tables, these will not be used because of the fact that I am instead using materialized views and this will greatly improve performance.
-- CREATE TABLE IF NOT EXISTS tblProductsCategoriesPrime (
--   ProductID INT NOT NULL,
--   CategoryID INT NOT NULL,
--   PRIMARY KEY (ProductID, CategoryID),
--   FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE,
--   FOREIGN KEY (CategoryID) REFERENCES tblCategoriesPrime (CategoryID) ON DELETE CASCADE
-- );

-- CREATE TABLE IF NOT EXISTS tblProductsCategoriesSub (
--   ProductID INT NOT NULL,
--   CategoryID INT NOT NULL,
--   PRIMARY KEY (ProductID, CategoryID),
--   FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE,
--   FOREIGN KEY (CategoryID) REFERENCES tblCategoriesSub (CategoryID) ON DELETE CASCADE
-- );

-- CREATE TABLE IF NOT EXISTS tblProductsCategoriesFinal (
--   ProductID INT NOT NULL,
--   CategoryID INT NOT NULL,
--   PRIMARY KEY (ProductID, CategoryID),
--   FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE,
--   FOREIGN KEY (CategoryID) REFERENCES tblCategoriesFinal (CategoryID) ON DELETE CASCADE
-- );


CREATE TABLE IF NOT EXISTS tblProductInventory (
  ProductID INT NOT NULL,
  Quantity INT NOT NULL,
  PRIMARY KEY (ProductID),
  FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblDiscount (
  DiscountID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  DiscountCode VARCHAR(255) NOT NULL,
  DiscountPercentage DECIMAL(5,2) NOT NULL,
  DiscountStartDate DATE,
  DiscountEndDate DATE
);

CREATE TABLE IF NOT EXISTS tblProductDiscount (
  ProductID INT NOT NULL,
  DiscountID INT NOT NULL,
  PRIMARY KEY (ProductID, DiscountID),
  FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE,
  FOREIGN KEY (DiscountID) REFERENCES tblDiscount (DiscountID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblImages (
  ImageID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ProductID INT NOT NULL,
  ImageURL VARCHAR(255) NOT NULL,
  FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblVariation (
  VariationID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  ProductID INT NOT NULL,
  VariationName VARCHAR(255) NOT NULL,
  VariationDescription TEXT,
  FOREIGN KEY (ProductID) REFERENCES tblProducts (ProductID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblAttribute (
  AttributeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  VariationID INT NOT NULL,
  AttributeName VARCHAR(255) NOT NULL,
  AttributeValue VARCHAR(255) NOT NULL,
  FOREIGN KEY (VariationID) REFERENCES tblVariation (VariationID) ON DELETE CASCADE
);

