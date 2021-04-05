
import json

# Write buyers into a dictionary data structure


def buyersToDictionary(filename):
    buyersFile = open(filename, "r+")
    data = json.load(buyersFile)
    dic_buyers = {}
    dic_buyers['buyers'] = {}
    for buyer in data['buyers']:
        dic_buyers['buyers'][buyer['id']] = {
            'buyer_name': buyer['name'],
            'buyer_age': buyer['age']
        }
    return dic_buyers


# get buyer detail by buyer_id
def getBuyerDetail(buyer_id):
    dic_buyers = buyersToDictionary("processed/buyers.json")
    # print(dic_buyers)
    # print(dic_buyers['buyers'][buyer_id])
    # print(buyer_id)
    return dic_buyers['buyers'][buyer_id]


"""
Read products in csv file and rewrite it in json format
"""

# Preprocess products into a object array data structure (JSON)


def productsToJson(filename):
    productsFile = open(filename, "r+")
    data = {}
    data['products'] = []
    no_line = 0
    for line in productsFile:
        # remove double quotes from product name
        if '\"' in line:
            line = line.replace("\"", "")
            for idx in range(len(line) - 1):
                letter = line[idx]
                # when finds simple quote among two lowercase characters, should be removed
                validation = (letter == '\'') and (
                    line[idx-1].isalpha() and line[idx-1].islower()) and (line[idx+1].isalpha() and line[idx+1].islower())
                if validation:
                    line = line[:idx] + line[idx+1:]
        line = line.strip().split('\'')
        # skip first line cause is headers
        if no_line > 0:
            # print(line)
            product_id = line[0]
            product_name = line[1]
            product_price = int(line[2])

            data['products'].append({
                'product_id': product_id,
                'product_name': product_name,
                'product_price': product_price
            })
        no_line += 1
    with open("files/products-processed-toObjectArray.json", "w+") as outfile:
        json.dump(data, outfile)
    outfile.close()


# Write products into a dictionary data structure
def productsToDict(filename):
    productsFile = open(filename, "r+")
    data = {}
    data['products'] = {}
    no_line = 0
    for line in productsFile:
        # remove double quotes from product name
        if '\"' in line:
            line = line.replace("\"", "")

            for idx in range(len(line) - 1):
                letter = line[idx]
                # when finds simple quote among two lowercase characters, should be removed
                validation = (letter == '\'') and (
                    line[idx-1].isalpha() and line[idx-1].islower()) and (line[idx+1].isalpha() and line[idx+1].islower())
                if validation:
                    line = line[:idx] + line[idx+1:]
        line = line.strip().split('\'')
        # skip first line cause is headers
        if no_line > 0:
            # print(line)
            product_id = line[0]
            product_name = line[1]
            product_price = int(line[2])

            data['products'][product_id] = {
                'product_name': product_name,
                'product_price': product_price
            }
        no_line += 1
    # with open("files/products-processed-toDictionary.json", "w+") as outfile:
    #     json.dump(data, outfile)
    # outfile.close()
    return data

# Get product detail by product_id


def getProductDetail(product_id):
    dic_products = productsToDict("raw/products.csv")
    return dic_products['products'][product_id]


"""
Read transactions in csv file and rewrite it in json format
"""


def getProductArray(arr_product_ids):
    arr = []
    for product_id in arr_product_ids:
        product_detail = getProductDetail(product_id)
        # print(product_detail)
        arr.append({
            'id': product_id,
            'name': product_detail['product_name'],
            'price': product_detail['product_price']
        })
    return arr


def preprocessTransactions(filename):
    transactionsFile = open(filename, "r+")
    data = {}
    data['transactions'] = []
    for line in transactionsFile:
        line = line.split('#')
        no_line = 0
        # print(line)
        for element in line:
            # skip first line cause is empty
            if no_line > 0:
                # print(element)
                # print(element.split(','))
                # print(element)

                transaction = element.split('_')
                # print(transaction)
                # break
                transaction_id = transaction[0]
                buyer_id = transaction[1]
                ip = transaction[2]
                device = transaction[3]

                arr_product_id = transaction[4]
                arr_product_id = arr_product_id.replace("(", "")
                arr_product_id = arr_product_id.replace(")", "")
                arr_product_id = arr_product_id.split(',')

                buyer_detail = getBuyerDetail(buyer_id)
                # print(buyer_detail)

                # for product_id in arr_product_id:
                #     data['transactions'].append({
                #         'transaction_id': transaction_id,
                #         'buyer': {
                #             'buyer_id': buyer_id,
                #             'buyer_name': buyer_detail['buyer_name'],
                #             'buyer_age': buyer_detail['buyer_age']
                #         },
                #         'product': [],
                #         'ip': ip,
                #         'device': device
                #     })

                data['transactions'].append({
                    'transaction_id': transaction_id,
                    # 'buyer': {
                    #     'id': buyer_id,
                    #     'name': buyer_detail['buyer_name'],
                    #     'age': buyer_detail['buyer_age']
                    # },
                    'buyer_id': buyer_id,
                    'buyer_name': buyer_detail['buyer_name'],
                    'buyer_age': buyer_detail['buyer_age'],
                    'products': getProductArray(arr_product_id),
                    'ip': ip,
                    'device': device,
                })
            no_line += 1
    with open("processed/transactions-processed.json", "w+") as outfile:
        json.dump(data, outfile)
    outfile.close()


# Call methods
# productsToJson("files/products.csv")
# buyersToDictionary("files/buyers.json")
preprocessTransactions("raw/transactions.csv")
