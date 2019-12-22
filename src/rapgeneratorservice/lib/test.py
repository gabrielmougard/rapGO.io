def square(x):
    newX = []
    for elt in x:
        newX.append(elt*elt)
    return newX

a = [1,2,3,4]
b = square(a)
print(b)