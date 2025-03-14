import numpy as np
import matplotlib.pyplot as plt
import pandas as pd
from sklearn.cluster import KMeans
from scipy import stats
from sklearn.mixture import GaussianMixture

def extract_data(data_file):

    df = pd.read_csv(data_file, header=None)

    boneid_to_points = {}
    total_points = []
    boneids = []

    print(df.shape)

    for _, row in df.iterrows():
        boneid = row[1]
        x_values = row[2::2].dropna().astype(float)  # Extract x values
        y_values = row[3::2].dropna().astype(float)

        points = list(zip(x_values, y_values))
        boneids.extend([boneid for i in range(len(points))])

        if boneid not in boneid_to_points:
            boneid_to_points[boneid] = []

        boneid_to_points[boneid].extend(points)
        total_points.extend(points)

    return boneid_to_points, total_points, boneids

def kmeans_plot(total_points, cluster_centers):
    # Extract keys (column index 1) and coordinate pairs (starting from index 2)
    plt.figure(figsize=(10, 6))
    x_values , y_values = zip(*total_points)

    plt.scatter(x_values, y_values, 
                s=200, c='blue', marker='o')


    plt.scatter(cluster_centers[:, 0], cluster_centers[:, 1], 
                s=200, c='red', marker='X')

    
    plt.title("Scatter Plot of (X, Y) Points Colored by Key")
    plt.legend()
    plt.grid(True)
    plt.show()

def gaussian_mixture_plot(points,labels, gmm, num_clusters = 4):
    plt.figure(figsize=(8, 6))
    plt.scatter(points[:, 0], points[:, 1], c=labels, cmap='viridis', alpha=0.7, label="Data Points")

    # Plot Gaussian ellipses for each cluster
    for i in range(num_clusters):
        mean = gmm.means_[i]
        cov = gmm.covariances_[i]

        # Get eigenvalues and eigenvectors for the covariance matrix
        eig_vals, eig_vecs = np.linalg.eigh(cov)
        angle = np.degrees(np.arctan2(*eig_vecs[:, 0][::-1]))  # Rotation angle
        width, height = 2 * np.sqrt(eig_vals)  # Scaling for 1 std deviation ellipse

        # Draw ellipse
        ellipse = plt.matplotlib.patches.Ellipse(mean, width, height, angle=angle, edgecolor='red', facecolor='none', lw=2)
        plt.gca().add_patch(ellipse)

    plt.xlabel("X Coordinate")
    plt.ylabel("Y Coordinate")
    plt.title(f"Gaussian Mixture Clustering (k={num_clusters})")
    plt.legend()
    plt.grid(True)
    plt.show()

def k_means_analysis(points, num_clusters = 4):

    kmeans = KMeans(n_clusters=num_clusters, random_state=42)
    print("Running kmeans...")
    labels = kmeans.fit_predict(points)
    print("Finished fitting kmeans!")

    return labels, kmeans.cluster_centers_

def guassian_mixture_clustering(points, num_clusters = 4):

    points = np.array(points)

    # Fit Gaussian Mixture Model
    gmm = GaussianMixture(n_components=num_clusters, covariance_type='full', random_state=42)
    labels = gmm.fit_predict(points)

    return points, labels, gmm

def create_groupings(labels, boneids):

    #I can implement the mode method

    groups = {0: [], 1: [], 2: [], 3: []}

    boneids = np.array(boneids)
    labels = np.array(labels)

    for i in range(10):
        idcs = np.where(boneids == i)[0]
        mode_result = stats.mode(labels[idcs])
        mode = mode_result.mode

        groups[mode].append(i)

    return groups

def main():

    boneid_to_points, points, boneids = extract_data("./bouncepoints.csv")
    # labels, cluster_centers = k_means_analysis(points, num_clusters=4)
    # groupings = create_groupings(labels, boneids)

    # print(groupings)

    # kmeans_plot(points, cluster_centers)
    points, labels, gmm = guassian_mixture_clustering(points, num_clusters = 4)

    print(labels)
    groupings = create_groupings(labels, boneids)

    print(groupings)
    gaussian_mixture_plot(points , labels, gmm)

if __name__ == "__main__":
    main()