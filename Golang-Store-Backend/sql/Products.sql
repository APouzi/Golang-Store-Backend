CREATE TABLE IF NOT EXISTS tblProducts (
  Product_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_Name VARCHAR(255) NOT NULL,
  Product_Description TEXT,
  Product_Price DECIMAL(10,2),
  SKU VARCHAR(50),
  UPC VARCHAR(50),
  PRIMARY_IMAGE VARCHAR(255) NULL,
  Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Modified_Date DATETIME NULL
  
);

CREATE TABLE IF NOT EXISTS tblProductVariation (
  Variation_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_ID INT NOT NULL,
  Variation_Name VARCHAR(255) NOT NULL,
  Variation_Description TEXT,
  Variation_Price DECIMAL(10,2) NOT NULL,
  SKU VARCHAR(50),
  UPC VARCHAR(50),
  PRIMARY_IMAGE VARCHAR(255) NULL,
  Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Modified_Date DATETIME NULL,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblProductInventory (
  Inv_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Variation_ID INT NOT NULL,
  Quantity INT NOT NULL,
  FOREIGN KEY (Variation_ID) REFERENCES tblProductVariation (Variation_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblLocation (
  Location_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Inv_ID INT NOT NULL,
  Location_At VARCHAR(255) NOT NULL,
  FOREIGN KEY (Inv_ID) REFERENCES tblProductInventory (Inv_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblCategoriesPrime (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesSub (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesFinal (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);

CREATE TABLE IF NOT EXISTS tblCatPrimeSub (
  CatPrimeID INT NOT NULL,
  CatSubID INT NOT NULL,
  FOREIGN KEY (CatPrimeID) REFERENCES tblCategoriesPrime (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (Category_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblCatSubFinal (
  CatSubID INT NOT NULL,
  CatFinalID INT NOT NULL,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (CatFinalID) REFERENCES tblCategoriesFinal (Category_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblCatFinalProd (
  CatFinalID INT NOT NULL,
  Product_ID INT NOT NULL,
  FOREIGN KEY (CatFinalID) REFERENCES tblCategoriesFinal (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
  
);


CREATE TABLE IF NOT EXISTS tblDiscount (
  Discount_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  DiscountCode VARCHAR(255) NOT NULL,
  DiscountPercentage DECIMAL(5,2) NOT NULL,
  DiscountStartDate DATE,
  DiscountEndDate DATE
);

CREATE TABLE IF NOT EXISTS tblProductDiscount (
  Product_ID INT NOT NULL,
  Discount_ID INT NOT NULL,
  PRIMARY KEY (Product_ID, Discount_ID),
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE,
  FOREIGN KEY (Discount_ID) REFERENCES tblDiscount (Discount_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblImages (
  ImageID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_ID INT NOT NULL,
  ImageURL VARCHAR(255) NOT NULL,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblVariation (
  Variation_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_ID INT NOT NULL,
  VariationName VARCHAR(255) NOT NULL,
  VariationDescription TEXT,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblAttribute (
  AttributeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Variation_ID INT NOT NULL,
  AttributeName VARCHAR(255) NOT NULL,
  AttributeValue VARCHAR(255) NOT NULL,
  FOREIGN KEY (Variation_ID) REFERENCES tblVariation (Variation_ID) ON DELETE CASCADE
);

CREATE VIEW PrimeSubFinalCategoryProducts AS SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblCategoriesPrime.CategoryName FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID ;